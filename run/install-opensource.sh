#!/usr/bin/env bash
set -euo pipefail

rapidou_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
requested_model="${RAPIDOU_QWEN_MODEL:-}"

if [[ $# -gt 2 ]]; then
    echo "usage: $0 [--model MODEL]" >&2
    exit 2
fi
if [[ $# -gt 0 ]]; then
    if [[ "$1" != "--model" || $# -ne 2 || -z "$2" ]]; then
        echo "usage: $0 [--model MODEL]" >&2
        exit 2
    fi
    requested_model="$2"
fi

need() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "error: required command is unavailable: $1" >&2
        exit 2
    }
}

api_ready() {
    curl -fsS --max-time 3 http://127.0.0.1:11434/api/tags >/dev/null 2>&1
}

wait_for_api() {
    local attempt
    for attempt in {1..30}; do
        api_ready && return 0
        sleep 1
    done
    echo "error: Ollama is installed but its API did not respond at http://127.0.0.1:11434" >&2
    exit 1
}

ollama_service_available() {
    command -v systemctl >/dev/null 2>&1 \
        && systemctl list-unit-files --type=service 2>/dev/null \
            | awk '{print $1}' \
            | grep -Fxq ollama.service
}

install_ollama() {
    local action="$1"
    if [[ "$(uname -s)" != "Linux" ]]; then
        echo "error: automatic Ollama installation and updates are Linux-only; install the latest version from https://ollama.com/download and rerun" >&2
        exit 2
    fi
    echo "$action Ollama with its official Linux installer..."
    curl -fsSL https://ollama.com/install.sh | sh
    hash -r
}

start_ollama() {
    if ollama_service_available; then
        if [[ "$(id -u)" -eq 0 ]]; then
            systemctl enable --now ollama
        else
            need sudo
            sudo systemctl enable --now ollama
        fi
    else
        ollama serve >"${TMPDIR:-/tmp}/rapidou-ollama.log" 2>&1 &
    fi
}

restart_ollama_after_update() {
    local pid executable arguments stopped=false
    if ollama_service_available; then
        if [[ "$(id -u)" -eq 0 ]]; then
            systemctl restart ollama
        else
            need sudo
            sudo systemctl restart ollama
        fi
        return
    fi

    while read -r pid executable arguments; do
        if [[ "${executable##*/}" == "ollama" && "$arguments" == "serve" ]]; then
            kill "$pid"
            stopped=true
        fi
    done < <(ps -eo pid=,args=)

    if [[ "$stopped" == true ]]; then
        for _ in {1..30}; do
            api_ready || break
            sleep 0.2
        done
    fi
    if ! api_ready; then
        start_ollama
    fi
}

pull_model() {
    local model="$1" pull_log
    pull_log="$(mktemp "${TMPDIR:-/tmp}/rapidou-ollama-pull.XXXXXX")"
    if ollama pull "$model" 2>&1 | tee "$pull_log"; then
        rm -f "$pull_log"
        return
    fi
    if ! grep -Fq "requires a newer version of Ollama" "$pull_log"; then
        rm -f "$pull_log"
        return 1
    fi

    rm -f "$pull_log"
    install_ollama "Updating"
    restart_ollama_after_update
    wait_for_api
    echo "Retrying $model with the updated Ollama..."
    ollama pull "$model"
}

detect_ram_gib() {
    local bytes kib
    if [[ -r /proc/meminfo ]]; then
        kib="$(awk '/^MemTotal:/ {print $2; exit}' /proc/meminfo)"
        echo $(( (kib + 1048575) / 1048576 ))
    elif command -v sysctl >/dev/null 2>&1; then
        bytes="$(sysctl -n hw.memsize)"
        echo $(( (bytes + 1073741823) / 1073741824 ))
    else
        echo "error: cannot detect total RAM; set RAPIDOU_QWEN_MODEL explicitly" >&2
        exit 2
    fi
}

select_model() {
    local ram="$1"
    if [[ -n "$requested_model" ]]; then
        printf '%s\n' "$requested_model"
    elif (( ram < 6 )); then
        echo "qwen3.5:0.8b"
    elif (( ram < 12 )); then
        echo "qwen3.5:2b"
    elif (( ram < 16 )); then
        echo "qwen3.5:4b"
    elif (( ram < 28 )); then
        echo "qwen3.5:9b"
    elif (( ram < 56 )); then
        echo "qwen3.8:27b"
    else
        echo "qwen3.8:27b-q8_0"
    fi
}

context_window_for() {
    case "$1" in
        qwen3.5:0.8b) echo 8192 ;;
        qwen3.5:2b|qwen3.5:4b) echo 16384 ;;
        *) echo 32768 ;;
    esac
}

report_gpu() {
    if command -v nvidia-smi >/dev/null 2>&1; then
        nvidia-smi --query-gpu=name,memory.total --format=csv,noheader 2>/dev/null | sed 's/^/  NVIDIA: /' || true
    elif command -v rocm-smi >/dev/null 2>&1; then
        rocm-smi --showproductname --showmeminfo vram 2>/dev/null | sed 's/^/  AMD: /' || true
    elif command -v lspci >/dev/null 2>&1; then
        local gpu
        gpu="$(lspci 2>/dev/null | grep -Ei 'VGA|3D controller' || true)"
        if [[ -n "$gpu" ]]; then
            printf '%s\n' "$gpu" | sed 's/^/  GPU: /'
        else
            echo "  GPU: not detected"
        fi
    else
        echo "  GPU: detection tools unavailable (RAM-only selection)"
    fi
}

need curl
need python3

if ! command -v ollama >/dev/null 2>&1; then
    install_ollama "Installing"
fi

if ! api_ready; then
    start_ollama
fi
wait_for_api

if ! command -v qwen >/dev/null 2>&1; then
    echo "Installing Qwen Code with its official standalone installer..."
    curl -fsSL https://qwen-code-assets.oss-cn-hangzhou.aliyuncs.com/installation/install-qwen-standalone.sh | bash -s -- --method standalone --source rapidou
    export PATH="$HOME/.local/bin:$HOME/bin:$PATH"
fi
need qwen

ram_gib="$(detect_ram_gib)"
model="$(select_model "$ram_gib")"
context_window="$(context_window_for "$model")"

echo
echo "Detected hardware:"
echo "  RAM: ${ram_gib} GiB"
report_gpu
echo "  Model: $model"
echo "  Context: $context_window tokens"

if ! ollama list | awk 'NR > 1 {print $1}' | grep -Fxq "$model"; then
    echo "Pulling $model..."
    pull_model "$model"
else
    echo "Model already present: $model"
fi

settings="$HOME/.qwen/settings.json"
python3 "$rapidou_root/ai/install/merge-qwen-settings.py" "$settings" "$model" "$context_window"

extension="qwenlm.qwen-code-vscode-ide-companion"
extension_status="not installed (VSCodium unavailable)"
if command -v codium >/dev/null 2>&1; then
    if codium --list-extensions 2>/dev/null | grep -Fxiq "$extension"; then
        extension_status="installed"
    elif codium --install-extension "$extension"; then
        extension_status="installed"
    else
        extension_status="manual installation required: $extension"
    fi
fi

"$rapidou_root/run/doctor-ai.sh" --smoke --model "$model"

echo
echo "Local Qwen setup complete."
echo
echo "Model:"
echo "  $model"
echo
echo "CLI:"
echo "  qwen"
echo
echo "Rapidou/headless:"
echo "  qwen -p \"Inspect this repository and summarize its harness.\" --output-format json --approval-mode plan"
echo
echo "VSCodium:"
echo "  $extension_status"
if [[ "$extension_status" == "installed" ]]; then
    echo
    echo "Remaining IDE step:"
    echo "  1. Restart/open VSCodium in the Rapidou repository"
    echo "  2. Open a new integrated terminal and run qwen"
    echo "  3. Run /ide enable"
elif [[ "$extension_status" == manual* ]]; then
    echo
    echo "Remaining IDE steps:"
    echo "  1. Install $extension from Open VSX"
    echo "  2. Restart/open VSCodium in the Rapidou repository"
    echo "  3. Open a new integrated terminal and run qwen"
    echo "  4. Run /ide enable"
fi
