# Skill index

- `craft`: complete prompt-to-tested-change workflow.
- `spec`: direct behavioral spec that doubles as the plan.
- `questions`: resolve material gaps one at a time.
- `select-agent-model`: mandatory Codex/Claude model selection before every subagent.
- `luna`: lower-complexity independent code search.
- `backend`: Go, SQLite, API behavior, and API functional tests.
- `frontend`: plain browser UI and click-driven Chromium tests.
- `run-functional-tests`: mandatory Docker/Chromium gate.
- `code-review`: independent spec and simplicity review.

Every delegated task must pass through `select-agent-model` and `ai/hooks/prepare-agent.sh`, whether it runs in parallel or sequentially. The installer exposes these canonical skills at `.agents/skills/` for Codex and `.claude/skills/` for Claude; do not maintain duplicate copies.
