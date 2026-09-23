#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 4 || "$1" != "select-agent-model" ]]; then
    echo "usage: $0 select-agent-model <codex|claude> <low|medium|high> <small|medium|large>" >&2
    echo "read ai/skills/select-agent-model/SKILL.md before invoking this hook" >&2
    exit 2
fi

platform="$2"
complexity="$3"
size="$4"

case "$platform" in
    codex)
        case "$complexity:$size" in
            low:small)
                model="gpt-5.6-terra"
                reasoning="low"
                ;;
            low:medium|medium:small)
                model="gpt-5.6-terra"
                reasoning="medium"
                ;;
            low:large|medium:medium)
                model="gpt-5.6-terra"
                reasoning="high"
                ;;
            medium:large)
                model="gpt-5.6-sol"
                reasoning="high"
                ;;
            high:small|high:medium)
                model="gpt-5.6-sol"
                reasoning="xhigh"
                ;;
            high:large)
                model="gpt-5.6-sol"
                reasoning="max"
                ;;
            *)
                echo "invalid classification: complexity=$complexity size=$size" >&2
                exit 2
                ;;
        esac
        ;;
    claude)
        case "$complexity:$size" in
            low:small|low:medium)
                model="haiku"
                ;;
            low:large|medium:small|medium:medium)
                model="sonnet"
                ;;
            medium:large|high:small|high:medium|high:large)
                model="opus"
                ;;
            *)
                echo "invalid classification: complexity=$complexity size=$size" >&2
                exit 2
                ;;
        esac
        ;;
    *)
        echo "invalid platform: $platform" >&2
        exit 2
        ;;
esac

if [[ -n "${RAPIDOU_SELECTION_DIR:-}" ]]; then
    state="$RAPIDOU_SELECTION_DIR/rapidou-agent-selection-$platform"
else
    repository="$(git rev-parse --show-toplevel 2>/dev/null || pwd -P)"
    repository_key="$(printf '%s' "$repository" | cksum | awk '{print $1}')"
    state="${TMPDIR:-/tmp}/rapidou-agent-selection-${UID:-user}-$repository_key-$platform"
fi

mkdir -p "$(dirname "$state")"
umask 077
{
    printf 'platform=%s\n' "$platform"
    printf 'model=%s\n' "$model"
    printf 'selected_at=%s\n' "$(date +%s)"
    if [[ "$platform" == "codex" ]]; then
        printf 'reasoning_effort=%s\n' "$reasoning"
    fi
} > "$state"

printf 'platform=%s\nmodel=%s\n' "$platform" "$model"
if [[ "$platform" == "codex" ]]; then
    printf 'reasoning_effort=%s\n' "$reasoning"
fi
