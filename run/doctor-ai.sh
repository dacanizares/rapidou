#!/usr/bin/env bash
set -euo pipefail

smoke=false
model="${RAPIDOU_QWEN_MODEL:-}"
while [[ $# -gt 0 ]]; do
    case "$1" in
        --smoke) smoke=true; shift ;;
        --model)
            [[ $# -ge 2 && -n "$2" ]] || { echo "error: --model requires a value" >&2; exit 2; }
            model="$2"
            shift 2
            ;;
        *) echo "usage: $0 [--smoke] [--model MODEL]" >&2; exit 2 ;;
    esac
done

status=0
ok() { printf '%-12s OK  %s\n' "$1" "$2"; }
bad() { printf '%-12s FAIL  %s\n' "$1" "$2"; status=1; }
note() { printf '%-12s --  %s\n' "$1" "$2"; }

if command -v ollama >/dev/null 2>&1; then
    ok "Ollama" "$(ollama --version 2>&1 | head -n 1)"
else
    bad "Ollama" "not installed"
fi

if command -v qwen >/dev/null 2>&1; then
    ok "Qwen Code" "$(qwen --version 2>&1 | head -n 1)"
else
    bad "Qwen Code" "not installed"
fi

if command -v curl >/dev/null 2>&1 && curl -fsS --max-time 3 http://127.0.0.1:11434/api/tags >/dev/null 2>&1; then
    ok "Ollama API" "http://127.0.0.1:11434"
else
    bad "Ollama API" "not responding"
fi

if [[ -z "$model" && -f "$HOME/.qwen/settings.json" ]] && command -v python3 >/dev/null 2>&1; then
    model="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1])).get("model", {}).get("name", ""))' "$HOME/.qwen/settings.json" 2>/dev/null || true)"
fi
if [[ -n "$model" ]] && command -v ollama >/dev/null 2>&1 && ollama list | awk 'NR > 1 {print $1}' | grep -Fxq "$model"; then
    ok "Model" "$model"
else
    bad "Model" "${model:-not configured}"
fi

extension="qwenlm.qwen-code-vscode-ide-companion"
if command -v codium >/dev/null 2>&1; then
    if codium --list-extensions 2>/dev/null | grep -Fxiq "$extension"; then
        ok "VSCodium" "Qwen Code Companion installed"
    else
        note "VSCodium" "manual step: install $extension"
    fi
else
    note "VSCodium" "not installed"
fi

if [[ "$smoke" == true && "$status" -eq 0 ]]; then
    temporary="$(mktemp -d)"
    trap 'rm -rf "$temporary"' EXIT
    output="$temporary/qwen-smoke.json"
    if (cd "$temporary" && qwen -p "Reply with RAPIDOU_OK. Do not use tools or modify files." \
        --output-format json --approval-mode plan --max-session-turns 2 --max-wall-time 5m >"$output") \
        && python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); assert isinstance(data,list); assert any(x.get("type")=="result" and x.get("subtype")=="success" for x in data if isinstance(x,dict))' "$output"; then
        ok "Headless" "JSON prompt completed in an isolated temp directory"
    else
        bad "Headless" "Qwen prompt failed"
    fi
fi

exit "$status"
