# Agent compatibility and installation

## Happy path

- Rapidou is pinned as a Git submodule at exactly `lib/rapidou`.
- From the consuming repository, `./lib/rapidou/run/install.sh` installs the shared workflow for all three clients without copying its source.
- Codex discovers every canonical skill through `.agents/skills/`; Claude discovers the same folders through `.claude/skills/`; Qwen discovers them through `.qwen/skills/`.
- The consuming repository keeps its own instructions. The installer appends a small managed Rapidou block to `AGENTS.md` and makes `CLAUDE.md` import it.
- After the client reloads the project and Codex trusts its repository hook, each client's `PreToolUse` hook blocks an independent agent until `prepare-agent.sh` records a fresh platform-specific selection. Codex and Claude require the selected model; Qwen inherits its active native model. One selection authorizes exactly one matching spawn.
- `./lib/rapidou/run/install-opensource.sh` installs and configures a local Qwen Code/Ollama stack without a cloud account, preserves unrelated user settings, and verifies a headless JSON prompt in a temporary directory.
- When `codium` is available, the local installer installs the official `qwenlm.qwen-code-vscode-ide-companion` extension. The user then restarts VSCodium, runs `qwen` in its integrated terminal, and enters `/ide enable` to connect the session.
- Running the installer again changes nothing and succeeds.
- `lib/rapidou/src/` remains the executable base example. It is reference code, not silently copied application source.

## Expected failures

- Installation fails when Rapidou is not located at `lib/rapidou`.
- Installation refuses to replace an existing skill or hook configuration it does not own. The error points to the fragment that must be merged manually.
- Unknown clients, complexity values, sizes, or mismatched spawn models are rejected.
- Invalid Qwen user settings are rejected before modification and an existing valid file is backed up once.
- A stale or already-consumed selection cannot authorize another agent.
- The consuming application must provide its own Docker/Chromium `./run/test.sh`; passing Rapidou's example suite alone never proves the consuming application works.

## Acceptance

- Automated harness tests cover both explicit model matrices, Qwen's inherited-model gate, installation, local-settings merge, idempotency, and safe conflict refusal.
- The README identifies Qwen Code Companion by its exact extension ID and documents both automatic and manual VSCodium installation.
- `./run/test.sh` passes with the real Chromium journey before and after independent review.
