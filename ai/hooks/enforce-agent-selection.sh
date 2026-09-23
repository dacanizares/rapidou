#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 || ( "$1" != "codex" && "$1" != "claude" && "$1" != "qwen" ) ]]; then
    echo "usage: $0 <codex|claude|qwen>" >&2
    exit 2
fi

expected_platform="$1"
payload="$(sed ':a;N;$!ba;s/[[:space:]]//g')"

if [[ -n "${RAPIDOU_SELECTION_DIR:-}" ]]; then
    state="$RAPIDOU_SELECTION_DIR/rapidou-agent-selection-$expected_platform"
else
    repository="$(git rev-parse --show-toplevel 2>/dev/null || pwd -P)"
    repository_key="$(printf '%s' "$repository" | cksum | awk '{print $1}')"
    state="${TMPDIR:-/tmp}/rapidou-agent-selection-${UID:-user}-$repository_key-$expected_platform"
fi

if [[ ! -f "$state" ]]; then
    echo "Rapidou blocked the subagent: run select-agent-model for $expected_platform immediately before spawning it." >&2
    exit 2
fi

claimed_state="$state.claimed.$$"
if ! mv "$state" "$claimed_state" 2>/dev/null; then
    echo "Rapidou blocked the subagent: this model selection was already consumed." >&2
    exit 2
fi
state="$claimed_state"
trap 'rm -f "$state"' EXIT

platform=""
model=""
reasoning_effort=""
selected_at=""
while IFS='=' read -r key value; do
    case "$key" in
        platform) platform="$value" ;;
        model) model="$value" ;;
        reasoning_effort) reasoning_effort="$value" ;;
        selected_at) selected_at="$value" ;;
    esac
done < "$state"

if [[ ! "$selected_at" =~ ^[0-9]+$ ]]; then
    echo "Rapidou blocked the subagent: invalid model selection state." >&2
    exit 2
fi
selection_age=$(( $(date +%s) - selected_at ))
if (( selection_age < 0 || selection_age > 300 )); then
    echo "Rapidou blocked the subagent: the model selection is stale; select again." >&2
    exit 2
fi
tool_input="${payload#*\"tool_input\":}"
if [[ "$tool_input" == "$payload" ]]; then
    echo "Rapidou blocked the subagent: hook input has no tool_input." >&2
    exit 2
fi
if [[ "$platform" != "$expected_platform" ]]; then
    echo "Rapidou blocked the subagent: selection was prepared for $platform, not $expected_platform." >&2
    exit 2
fi
if [[ "$expected_platform" != "qwen" && "$tool_input" != *"\"model\":\"$model\""* ]]; then
    echo "Rapidou blocked the subagent: use the exact selected model $model for $expected_platform." >&2
    exit 2
fi
if [[ "$expected_platform" == "codex" && "$tool_input" != *"\"reasoning_effort\":\"$reasoning_effort\""* ]]; then
    echo "Rapidou blocked the subagent: use reasoning_effort=$reasoning_effort." >&2
    exit 2
fi
