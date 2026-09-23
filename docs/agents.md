# Coding agent compatibility

Verified against current official documentation on 2026-09-23. Rapidou does not provide an agent loop. Claude Code, Codex, and Qwen Code consume the same repository contract and keep their native execution, permissions, MCP, and subagent behavior.

## Portable layer

| Capability | Rapidou contract | Claude Code | Codex | Qwen Code |
| --- | --- | --- | --- | --- |
| Project instructions | Root `AGENTS.md`, with nested files when an application needs narrower rules | Reads `AGENTS.md` directly in current versions; the tiny `CLAUDE.md` import remains compatible | Reads root-to-working-directory `AGENTS.md` files | Reads `AGENTS.md` directly |
| Skills | Canonical folders under `ai/skills/` | `.claude/skills/` links | `.agents/skills/` links | `.qwen/skills/` links |
| Shell and file work | Native client tools execute the repository's plain scripts and edit normal files | Native tools | Native tools | Native tools |
| Verification | Repository-owned `run/` scripts and exit status | Supported | Supported | Supported |
| Headless use | Client-native non-interactive command | `claude -p ... --output-format json` | `codex exec --json ...` | `qwen -p ... --output-format json` |

All three clients support hooks, MCP servers, skills, shell execution, file editing, subagents, permission controls, and structured/headless output. These are compatible capabilities, not identical configuration schemas.

Instruction traversal stays client-native: Codex concatenates one instruction file per directory from the project root to the working directory; Claude loads ancestor instructions at launch and discovers narrower files as it enters subdirectories; Qwen loads its hierarchical context and existing `AGENTS.md` compatibility. Rapidou relies only on the shared root contract and does not require identical precedence rules. Each headless command preserves a process exit status suitable for scripts; its JSON envelope remains client-specific.

## Thin adapters

- `CLAUDE.md` imports `@AGENTS.md`; this remains useful for older Claude Code releases even though Claude Code 2.1.277 and later can read `AGENTS.md` directly.
- `run/install.sh` links the canonical skills into each client's discovery directory and installs its native hook file. It never duplicates skill bodies.
- Codex and Claude model selection maps complexity and size to an explicit client model. Qwen's native `agent` tool has no per-call model field, so Rapidou records `model=inherit`; Qwen keeps its active model or a named subagent's native model configuration.
- `run/install-opensource.sh` configures only the user-local Qwen/Ollama provider. Qwen-specific settings do not enter application source.

## Intentionally native

Do not try to unify approval modes, sandbox implementations, MCP configuration files, subagent schemas, event payloads, or structured-output envelopes. The common contract is observable repository behavior: load the instructions, invoke canonical skills, use the scripts, preserve exit status, and satisfy the application's verification gate.

Codex uses its OS sandbox and approval policy; Claude Code uses its permissions and optional sandbox; Qwen Code provides `plan`, `default`, `auto-edit`, `auto`, and `yolo` approval modes plus its own sandbox. Automation should choose the least privilege each client supports rather than translating modes by name.

## Local Qwen path

Run from the Rapidou checkout:

```sh
./run/install-opensource.sh
```

The script installs missing Ollama and Qwen Code on Linux using their official installers, starts Ollama, chooses a conservative model from total RAM, pulls it once, merges an OpenAI-compatible provider into `~/.qwen/settings.json`, installs the official Qwen Code Companion in VSCodium when available, and runs an isolated headless smoke prompt. Override only the model when needed:

```sh
RAPIDOU_QWEN_MODEL=qwen3.5:9b ./run/install-opensource.sh
```

or:

```sh
./run/install-opensource.sh --model qwen3.5:9b
```

The merge preserves unrelated settings and providers. Before its first modification of an existing file it creates `~/.qwen/settings.json.rapidou-backup`. The local endpoint is `http://127.0.0.1:11434/v1`; the placeholder key is stored under Rapidou's own environment name because Ollama requires a value but ignores it.

Run the lightweight diagnosis at any time:

```sh
./run/doctor-ai.sh
./run/doctor-ai.sh --smoke
```

Qwen Code Companion is published on Open VSX as `qwenlm.qwen-code-vscode-ide-companion`. The installer adds it automatically when `codium` is on `PATH`. Install it manually when necessary with:

```sh
codium --install-extension qwenlm.qwen-code-vscode-ide-companion
```

It can also be found from VSCodium's Extensions view by searching for that exact ID. After installation, restart VSCodium, open a new integrated terminal in the application repository, start `qwen`, and run `/ide enable`; that connection step requires the IDE environment and remains interactive. Microsoft VS Code accepts the same extension ID through `code --install-extension`, but Rapidou only automates VSCodium installation.

## Verified official references

- Claude Code: [project instructions](https://code.claude.com/docs/en/memory), [skills](https://code.claude.com/docs/en/skills), [hooks](https://code.claude.com/docs/en/hooks), [subagents](https://code.claude.com/docs/en/sub-agents), [MCP](https://code.claude.com/docs/en/mcp), [CLI/headless flags](https://code.claude.com/docs/en/cli-reference), and [permissions](https://code.claude.com/docs/en/permissions).
- Codex: [AGENTS.md](https://learn.chatgpt.com/docs/agent-configuration/agents-md), [skills](https://learn.chatgpt.com/docs/skills), [hooks](https://learn.chatgpt.com/docs/config-advanced#hooks), [MCP](https://learn.chatgpt.com/docs/extend/mcp), [subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents), [non-interactive mode](https://learn.chatgpt.com/docs/non-interactive-mode), and [sandbox/approvals](https://learn.chatgpt.com/docs/sandboxing).
- Qwen Code: [AGENTS.md and memory](https://qwenlm.github.io/qwen-code-docs/en/users/features/memory/), [skills](https://qwenlm.github.io/qwen-code-docs/en/users/features/skills/), [hooks](https://qwenlm.github.io/qwen-code-docs/en/users/features/hooks/), [agent tool](https://qwenlm.github.io/qwen-code-docs/en/developers/tools/task/), [MCP](https://qwenlm.github.io/qwen-code-docs/en/users/features/mcp/), [headless mode](https://qwenlm.github.io/qwen-code-docs/en/users/features/headless/), [model providers](https://qwenlm.github.io/qwen-code-docs/en/users/configuration/model-providers/), and [IDE integration](https://qwenlm.github.io/qwen-code-docs/en/users/ide-integration/ide-integration/).
- Local runtime: [Ollama Linux installation](https://docs.ollama.com/linux), [OpenAI compatibility](https://docs.ollama.com/api/openai-compatibility), [Qwen3.5 tags](https://ollama.com/library/qwen3.5/tags), [Qwen3.8 tags](https://ollama.com/library/qwen3.8/tags), and [Qwen Code installation](https://qwenlm.github.io/qwen-code-docs/en/users/overview/).
