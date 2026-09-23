#!/usr/bin/env bash
set -euo pipefail

rapidou_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
prompt_file="$rapidou_root/ai/benchmarks/instagram-like-prompt.md"
output_root=""
selection="all"
turns=80
wall_time=90m
prepare=false
failed_trials=0

usage() {
    cat <<'EOF'
usage: run/benchmark-instagram.sh --output DIRECTORY [options]

Runs comparable isolated Instagram-like application builds.

  --output DIRECTORY     New directory that will receive one Git worktree per trial.
  --run NAME             all (default), codex, codex-oss-qwen27b, qwen-9b, or qwen-27b.
  --prepare              Pull both local Ollama models and exit; setup time is not measured.
  --turns NUMBER         Maximum agent turns per Qwen trial (default: 80).
  --wall-time DURATION   Maximum wall time per trial (default: 90m).

`codex` uses the authenticated Codex default model. `codex-oss-qwen27b` uses
Codex's local Ollama provider and Qwen 3.8 27B. The latter measures Codex's
agent/tool loop with a local Qwen model; it is not cloud-Codex subagent
orchestration.
EOF
}

need() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "error: required command is unavailable: $1" >&2
        exit 2
    }
}

model_present() {
    ollama list | awk 'NR > 1 {print $1}' | grep -Fxq "$1"
}

prepare_models() {
    local model
    local models=()
    need ollama
    case "$selection" in
        all) models=(qwen3.5:9b qwen3.8:27b) ;;
        qwen-9b) models=(qwen3.5:9b) ;;
        qwen-27b|codex-oss-qwen27b) models=(qwen3.8:27b) ;;
        codex) return ;;
    esac
    for model in "${models[@]}"; do
        if model_present "$model"; then
            echo "Model already present: $model"
        else
            echo "Pulling $model (benchmark setup; not timed)..."
            ollama pull "$model"
        fi
    done
}

write_metadata() {
    local destination="$1" trial="$2" runner="$3" model="$4" started="$5" finished="$6" status="$7"
    python3 - "$destination" "$trial" "$runner" "$model" "$started" "$finished" "$status" "$prompt_file" <<'PY'
import hashlib
import json
import pathlib
import sys

destination, trial, runner, model, started, finished, status, prompt_path = sys.argv[1:]
prompt = pathlib.Path(prompt_path).read_bytes()
data = {
    "trial": trial,
    "runner": runner,
    "model": model,
    "started_at_unix": int(started),
    "finished_at_unix": int(finished),
    "elapsed_seconds": int(finished) - int(started),
    "exit_status": int(status),
    "prompt_sha256": hashlib.sha256(prompt).hexdigest(),
}
pathlib.Path(destination).write_text(json.dumps(data, indent=2) + "\n")
PY
}

run_trial() {
    local trial="$1" runner="$2" model="$3"
    local trial_root="$output_root/$trial" worktree="$trial_root/app"
    local started finished status

    mkdir -p "$trial_root"
    git -C "$rapidou_root" worktree add --detach "$worktree" HEAD >/dev/null
    cp "$prompt_file" "$trial_root/prompt.md"

    echo
    echo "Starting $trial in $worktree"
    started="$(date +%s)"
    set +e
    if [[ "$runner" == "qwen" ]]; then
        (
            cd "$worktree"
            qwen -p "$(<"$prompt_file")" \
                --auth-type openai \
                --openai-api-key ollama \
                --openai-base-url http://127.0.0.1:11434/v1 \
                --model "$model" \
                --output-format json \
                --approval-mode auto-edit \
                --max-session-turns "$turns" \
                --max-wall-time "$wall_time"
        ) >"$trial_root/agent.json" 2>&1
        status=$?
    elif [[ "$runner" == "codex-oss" ]]; then
        (
            cd "$worktree"
            codex exec \
                --oss \
                --local-provider ollama \
                --model "$model" \
                --sandbox workspace-write \
                --json \
                --ephemeral \
                "$(<"$prompt_file")"
        ) >"$trial_root/agent.jsonl" 2>&1
        status=$?
    else
        (
            cd "$worktree"
            codex exec \
                --sandbox workspace-write \
                --json \
                --ephemeral \
                "$(<"$prompt_file")"
        ) >"$trial_root/agent.jsonl" 2>&1
        status=$?
    fi
    set -e
    finished="$(date +%s)"
    write_metadata "$trial_root/metadata.json" "$trial" "$runner" "$model" "$started" "$finished" "$status"
    printf '%s: %ss, exit %s\n' "$trial" "$(( finished - started ))" "$status"
    if (( status != 0 )); then
        failed_trials=1
    fi
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --output)
            [[ $# -ge 2 && -n "$2" ]] || { usage >&2; exit 2; }
            output_root="$2"
            shift 2
            ;;
        --run)
            [[ $# -ge 2 ]] || { usage >&2; exit 2; }
            selection="$2"
            shift 2
            ;;
        --prepare) prepare=true; shift ;;
        --turns)
            [[ $# -ge 2 && "$2" =~ ^[1-9][0-9]*$ ]] || { usage >&2; exit 2; }
            turns="$2"
            shift 2
            ;;
        --wall-time)
            [[ $# -ge 2 && -n "$2" ]] || { usage >&2; exit 2; }
            wall_time="$2"
            shift 2
            ;;
        --help|-h) usage; exit 0 ;;
        *) usage >&2; exit 2 ;;
    esac
done

case "$selection" in
    all|codex|codex-oss-qwen27b|qwen-9b|qwen-27b) ;;
    *) echo "error: unknown trial: $selection" >&2; exit 2 ;;
esac

if [[ "$prepare" == false && -z "$output_root" ]]; then
    echo "error: --output is required unless --prepare is used" >&2
    exit 2
fi
if [[ "$prepare" == true ]]; then
    prepare_models
    exit 0
fi
[[ ! -e "$output_root" ]] || { echo "error: output directory already exists: $output_root" >&2; exit 2; }
need git
need python3
if [[ "$selection" == all || "$selection" == qwen-9b || "$selection" == qwen-27b || "$selection" == codex-oss-qwen27b ]]; then
    need qwen
fi
if [[ "$selection" == all || "$selection" == codex || "$selection" == codex-oss-qwen27b ]]; then
    need codex
fi

prepare_models

mkdir -p "$output_root"
if [[ "$selection" == all || "$selection" == codex ]]; then
    run_trial codex codex codex-default
fi
if [[ "$selection" == all || "$selection" == codex-oss-qwen27b ]]; then
    run_trial codex-oss-qwen27b codex-oss qwen3.8:27b
fi
if [[ "$selection" == all || "$selection" == qwen-9b ]]; then
    run_trial qwen-9b qwen qwen3.5:9b
fi
if [[ "$selection" == all || "$selection" == qwen-27b ]]; then
    run_trial qwen-27b qwen qwen3.8:27b
fi

echo "\nBenchmark results: $output_root"
exit "$failed_trials"
