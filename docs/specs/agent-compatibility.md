# Agent compatibility and installation

## Happy path

- Rapidou is pinned as a Git submodule at exactly `lib/rapidou`.
- From the consuming repository, `./lib/rapidou/run/install` installs the shared workflow for both clients without copying its source.
- Codex discovers every canonical skill through `.agents/skills/`; Claude discovers the same folders through `.claude/skills/`.
- The consuming repository keeps its own instructions. The installer appends a small managed Rapidou block to `AGENTS.md` and makes `CLAUDE.md` import it.
- After the client reloads the project and Codex trusts its repository hook, each client's `PreToolUse` hook blocks an independent agent until `prepare-agent.sh` records a fresh platform-specific model selection. One selection authorizes exactly one matching spawn.
- Running the installer again changes nothing and succeeds.
- `lib/rapidou/src/` remains the executable base example. It is reference code, not silently copied application source.

## Expected failures

- Installation fails when Rapidou is not located at `lib/rapidou`.
- Installation refuses to replace an existing skill or hook configuration it does not own. The error points to the fragment that must be merged manually.
- Unknown clients, complexity values, sizes, or mismatched spawn models are rejected.
- A stale or already-consumed selection cannot authorize another agent.
- The consuming application must provide its own Docker/Chromium `./run/test`; passing Rapidou's example suite alone never proves the consuming application works.

## Acceptance

- Automated tests cover both model matrices, one-shot hook enforcement, installation, idempotency, and safe conflict refusal.
- `./run/test` passes with the real Chromium journey before and after independent review.
