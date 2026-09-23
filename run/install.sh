#!/usr/bin/env bash
set -euo pipefail

rapidou_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

accept_existing_hooks=false
if [[ "${1:-}" == "--accept-existing-hooks" ]]; then
    accept_existing_hooks=true
    shift
fi

if [[ $# -gt 1 ]]; then
    echo "usage: $0 [--accept-existing-hooks] [project-root]" >&2
    exit 2
fi

if [[ $# -eq 1 ]]; then
    project_root="$1"
else
    project_root="$(git -C "$rapidou_root" rev-parse --show-superproject-working-tree 2>/dev/null || true)"
    if [[ -z "$project_root" ]]; then
        echo "error: Rapidou is not a submodule; pass the consuming project root explicitly" >&2
        exit 2
    fi
fi

project_root="$(cd "$project_root" && pwd -P)"
expected_root="$project_root/lib/rapidou"
if [[ ! -d "$expected_root" || ! "$expected_root" -ef "$rapidou_root" ]]; then
    echo "error: Rapidou must be installed at $project_root/lib/rapidou" >&2
    exit 2
fi

skills=(backend code-review craft frontend luna questions run-functional-tests select-agent-model spec)

check_link() {
    local destination="$1"
    local expected="$2"
    if [[ -L "$destination" ]]; then
        if [[ "$(readlink "$destination")" != "$expected" ]]; then
            echo "error: existing link has a different target: $destination" >&2
            exit 2
        fi
    elif [[ -e "$destination" ]]; then
        echo "error: refusing to replace existing path: $destination" >&2
        exit 2
    fi
}

check_config() {
    local destination="$1"
    local expected="$2"
    local platform="$3"
    if [[ -L "$destination" ]]; then
        if [[ "$(readlink "$destination")" != "$expected" ]]; then
            echo "error: existing configuration link has a different target: $destination" >&2
            exit 2
        fi
    elif [[ -f "$destination" ]]; then
        if [[ "$accept_existing_hooks" != true ]]; then
            echo "error: $destination already exists; merge the Rapidou PreToolUse entry, then rerun with --accept-existing-hooks" >&2
            exit 2
        fi
        for required in '"PreToolUse"' '"matcher": "Agent"' 'lib/rapidou/ai/hooks/enforce-agent-selection.sh' "$platform"; do
            if ! grep -Fq "$required" "$destination"; then
                echo "error: existing configuration is missing Rapidou hook data: $destination" >&2
                exit 2
            fi
        done
    elif [[ -e "$destination" ]]; then
        echo "error: refusing to replace existing path: $destination" >&2
        exit 2
    fi
}

for skill in "${skills[@]}"; do
    check_link "$project_root/.agents/skills/$skill" "../../lib/rapidou/ai/skills/$skill"
    check_link "$project_root/.claude/skills/$skill" "../../lib/rapidou/ai/skills/$skill"
done
check_config "$project_root/.codex/hooks.json" "../lib/rapidou/ai/install/codex-hooks.json" codex
check_config "$project_root/.claude/settings.json" "../lib/rapidou/ai/install/claude-settings.json" claude

agents_start='<!-- rapidou:start -->'
agents_end='<!-- rapidou:end -->'
if [[ -f "$project_root/AGENTS.md" ]] && grep -Fq "$agents_start" "$project_root/AGENTS.md" && ! grep -Fq "$agents_end" "$project_root/AGENTS.md"; then
    echo "error: incomplete Rapidou block in $project_root/AGENTS.md" >&2
    exit 2
fi
claude_start='<!-- rapidou-claude:start -->'
claude_end='<!-- rapidou-claude:end -->'
if [[ -f "$project_root/CLAUDE.md" ]] && grep -Fq "$claude_start" "$project_root/CLAUDE.md" && ! grep -Fq "$claude_end" "$project_root/CLAUDE.md"; then
    echo "error: incomplete Rapidou block in $project_root/CLAUDE.md" >&2
    exit 2
fi

mkdir -p "$project_root/.agents/skills" "$project_root/.claude/skills" "$project_root/.codex"
for skill in "${skills[@]}"; do
    [[ -L "$project_root/.agents/skills/$skill" ]] || ln -s "../../lib/rapidou/ai/skills/$skill" "$project_root/.agents/skills/$skill"
    [[ -L "$project_root/.claude/skills/$skill" ]] || ln -s "../../lib/rapidou/ai/skills/$skill" "$project_root/.claude/skills/$skill"
done
[[ -e "$project_root/.codex/hooks.json" ]] || ln -s "../lib/rapidou/ai/install/codex-hooks.json" "$project_root/.codex/hooks.json"
[[ -e "$project_root/.claude/settings.json" ]] || ln -s "../lib/rapidou/ai/install/claude-settings.json" "$project_root/.claude/settings.json"

if [[ -f "$project_root/AGENTS.md" ]] && grep -Fq "$agents_start" "$project_root/AGENTS.md"; then
    :
else
    {
        printf '\n%s\n' "$agents_start"
        printf '## Rapidou\n\n'
        printf 'Rapidou is installed at `lib/rapidou`. Read `lib/rapidou/docs/index.md` and use the `craft` skill for implementation work.\n\n'
        printf 'Before every independent subagent, read `select-agent-model` and run `./lib/rapidou/ai/hooks/prepare-agent.sh select-agent-model <codex|claude> <complexity> <size>`. Use the exact returned model and, for Codex, reasoning effort.\n\n'
        printf 'The consuming application must own its Docker/Chromium `./run/test.sh` completion gate. `lib/rapidou/src/` is the executable base example, not application-owned source.\n'
        printf '%s\n' "$agents_end"
    } >> "$project_root/AGENTS.md"
fi

if [[ -f "$project_root/CLAUDE.md" ]] && grep -Fq "$claude_start" "$project_root/CLAUDE.md"; then
    :
else
    {
        printf '\n%s\n' "$claude_start"
        printf '@AGENTS.md\n'
        printf '%s\n' "$claude_end"
    } >> "$project_root/CLAUDE.md"
fi

printf 'Rapidou installed for Codex and Claude in %s\n' "$project_root"
printf 'Restart both clients; in Codex, open /hooks and trust the project hook.\n'
printf 'Next: adapt lib/rapidou/src/ as the base example and make the application own ./run/test.sh.\n'
