# Harness workflow

The files under `tst/` have separate roles:

- `shared.go`: shared HTTP types and helpers used by the functional journeys;
- `api.go`: application API and seed journeys;
- `ui.go`: the real Chromium click journey;
- `workflow_test.go`: self-tests for Rapidou's skills, agent hooks, installer, and documentation routing. It tests the harness workflow itself, not Museo Pixel behavior.

## Happy path

- `craft` receives a feature or bug-fix prompt and uses `spec` to update a concise observable-behavior contract; that spec is the plan.
- The spec names the happy path, realistic user mistakes, expected errors, and required functional coverage.
- `questions` resolves only material gaps, one question at a time, before implementation.
- Before every independent agent, the parent invokes `select-agent-model`, classifies complexity and size, and runs `ai/hooks/prepare-agent.sh` for Codex or Claude.
- The selector returns the exact platform model and, for Codex, reasoning effort. This applies to parallel and sequential agents, Luna searches, and reviews.
- Codex and Claude `PreToolUse` hooks intercept their `Agent` tool and consume one matching selection. `AGENTS.md` remains the readable policy.
- Backend and frontend agents read their domain indexes, implement only their owned scope, and add the corresponding API or click-driven browser tests.
- `run-functional-tests` runs the complete Docker/Chromium gate after implementation.
- A fresh independent `code-review` agent reviews the spec, diff, tests, and simplicity rules without editing.
- Valid review findings are fixed and the complete Docker/Chromium gate runs again on the final tree.

## Workflow failures

- Hook invocation without the literal `select-agent-model` argument is rejected, and direct hook bypass violates the mandatory `AGENTS.md` delegation gate.
- A spawn without a fresh selection, with the wrong platform, or with a different model is blocked.
- An unknown complexity or size classification is rejected.
- A container engine other than Docker or Podman is rejected.
- An unanswered material product question pauses implementation; it is never silently decided by an agent.
- Overlapping backend/frontend file ownership is serialized or reassigned before spawning agents.
- Missing, skipped, or failing API, Chromium UI, seed, or delegation workflow results fail `./run/test`.
- Host-only tests, API-only tests, a successful image build, or manual HTTP checks never satisfy completion.
