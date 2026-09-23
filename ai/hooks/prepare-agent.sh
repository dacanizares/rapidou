#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 3 || "$1" != "select-agent-model" ]]; then
    echo "usage: $0 select-agent-model <low|medium|high> <small|medium|large>" >&2
    echo "read ai/skills/select-agent-model/SKILL.md before invoking this hook" >&2
    exit 2
fi

complexity="$2"
size="$3"

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

printf 'model=%s\nreasoning_effort=%s\n' "$model" "$reasoning"
