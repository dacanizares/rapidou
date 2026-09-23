We want Rapidou to support **Claude Code, OpenAI Codex, and Qwen Code as first-class coding agents**.

The goal is NOT to build our own agent loop and NOT to normalize every agent down to the lowest common denominator.

Rapidou should provide the harness, project conventions, skills, tests, scripts, and verification. Claude Code, Codex, and Qwen Code should each be able to operate on the same Rapidou application using their native capabilities.

A user with no paid AI subscription must be able to use Rapidou entirely locally using:

* Qwen Code
* a local Qwen model
* Ollama
* optionally the Qwen Code Companion extension in VSCodium

I want this setup to be as direct and boring as possible.

## 1. Inspect the repository first

Before changing anything:

* inspect the current Rapidou repository structure
* inspect `AGENTS.md`, `CLAUDE.md`, `ai/`, `ai/skills/`, `ai/hooks/`, `run/`, `doc/`, tests, and any existing agent-specific configuration
* understand how Rapidou currently expects an AI agent to build, test, and verify an application
* do not duplicate existing mechanisms
* do not introduce a framework or unnecessary abstraction layer

Rapidou deliberately prefers:

* simple files
* shell scripts
* explicit conventions
* minimal dependencies
* no ceremony
* easy debugging
* predictable behavior

## 2. Verify Claude Code / Codex / Qwen Code compatibility

Research the CURRENT official documentation for all three tools before implementing compatibility assumptions.

At minimum verify:

* project instruction files
* nested/project instructions
* skills
* hooks
* MCP
* subagents
* shell execution
* file editing
* non-interactive/headless execution
* structured output, if available
* sandbox / permission modes
* exit codes suitable for automation
* ability to use Rapidou's existing scripts and tests

Create a concise compatibility document under `doc/`.

Do not write a generic feature comparison.

The important question is:

> Can Claude Code, Codex, and Qwen Code all correctly consume and operate the Rapidou harness?

Identify:

1. the common portable layer
2. thin agent-specific adapters that are actually required
3. features that cannot or should not be unified

Do NOT hide useful native capabilities just to make the three tools look identical.

## 3. Prefer one canonical project instruction source

Investigate whether `AGENTS.md` can reasonably be the canonical agent-neutral Rapidou instruction file.

Current expected direction, VERIFY before relying on it:

* Codex reads `AGENTS.md`
* Qwen Code reads `AGENTS.md`
* Claude Code can have a very small `CLAUDE.md` importing `@AGENTS.md`

If this is still supported by the current tools, use that pattern.

Avoid maintaining three copies of the same project instructions.

Agent-specific instruction files should contain only the minimum adapter instructions genuinely needed by that agent.

Likewise, inspect whether existing Rapidou skills can remain in one canonical location and be exposed to each agent using thin adapters/symlinks/imports where supported.

Do not duplicate skill bodies unless there is no reliable alternative.

## 4. Add the local/open-source installation path

Create:

```
run/install-opensource.sh
```

The intended user experience is:

```
./run/install-opensource.sh
```

After that command, the user should have the local Qwen coding stack installed and configured as far as can safely be automated.

The script must be:

* idempotent
* readable
* reasonably short
* safe to run repeatedly
* fail-fast with useful errors
* non-destructive to unrelated user configuration
* Linux-first
* Fedora must work well
* Debian/Ubuntu support where straightforward
* macOS support is welcome if it does not complicate the script significantly

Do not introduce a configuration-management framework.

Plain shell is preferred.

## 5. Install Ollama

If Ollama is missing, install it using the CURRENT official Ollama installation method.

On Linux the current official installer is expected to be similar to:

```
curl -fsSL https://ollama.com/install.sh | sh
```

VERIFY the current command.

Make sure Ollama is running before continuing.

Where systemd is available, use it appropriately.

Do not silently break an existing Ollama installation.

Detect and report:

* Ollama version
* whether its local API responds
* relevant GPU information when easily available

AMD/Linux must not be assumed to be NVIDIA-only.

## 6. Install Qwen Code

If `qwen` is missing, install Qwen Code using its CURRENT official standalone installer.

Prefer the standalone installer instead of introducing Node/npm as a requirement merely to install Qwen Code.

VERIFY the current official installation command before committing it.

After installation verify:

```
qwen --version
```

We need the CLI even if the editor extension also bundles Qwen Code because Rapidou must be able to execute Qwen headlessly from scripts.

## 7. Select a sensible local model

The installer must choose a local model conservatively based on available hardware.

Keep the logic SIMPLE.

RAM is the required baseline signal.

GPU/VRAM can improve the selection when reliable detection is straightforward, but do not build a hardware-detection framework.

Current rough desired behavior:

```
~4 GB       -> very small Qwen model
~8 GB       -> small Qwen model
~16 GB      -> medium local Qwen
~32 GB      -> strong coding Qwen
~64-128 GB  -> strongest sensible local Qwen configuration
```

Current candidate family to investigate:

```
qwen3.5:0.8b
qwen3.5:2b
qwen3.5:4b
qwen3.5:9b
qwen3.8:27b
qwen3.8:27b-q8_0
```

DO NOT blindly use these tags.

Check the CURRENT Ollama registry/documentation and select the best currently available Qwen coding/agentic models.

Prefer Qwen3.8 where it materially improves coding/agent behavior and fits safely.

The selection must leave enough RAM for:

* OS
* Ollama
* context/KV cache
* Qwen Code
* build/test processes

Do not choose a model merely because its model file technically fits into total RAM.

A 4 GB machine is expected to run only a very small model and will naturally require smaller tasks.

Allow explicit override, preferably something simple such as:

```
RAPIDOU_QWEN_MODEL=qwen3.5:9b ./run/install-opensource.sh
```

or a similarly simple `--model` argument.

Do not add ten command-line flags.

Print the detected hardware and selected model clearly.

Then:

```
ollama pull <selected-model>
```

unless the model already exists locally.

## 8. Configure Qwen Code to use local Ollama

Configure Qwen Code to talk to the local Ollama OpenAI-compatible endpoint.

Current expected endpoint:

```
http://127.0.0.1:11434/v1
```

VERIFY against current Qwen Code and Ollama documentation.

Configure the selected model using the current supported Qwen Code `modelProviders` mechanism.

Do not use deprecated Qwen configuration fields.

Do not require any cloud account or API key.

If Qwen requires a placeholder API key for an unauthenticated OpenAI-compatible local endpoint, configure the documented placeholder such as:

```
OLLAMA_API_KEY=ollama
```

Do not overwrite an existing `~/.qwen/settings.json`.

Merge only the configuration Rapidou needs, preserving unrelated providers and settings.

If modifying a user file, make a sensible backup before the first modification.

Do not commit machine-specific model choices or user credentials to the Rapidou repository.

## 9. Verify headless Qwen

Rapidou needs Qwen Code as an executable agent, not merely chat.

Verify the CURRENT equivalent of:

```
qwen -p "..." --output-format json
```

and/or `stream-json`.

Add a very small smoke test that verifies:

* Qwen Code starts
* it talks to the local Ollama model
* a simple headless prompt succeeds
* exit status can be consumed by another script

The smoke test must not modify the Rapidou repository.

Use a temp directory if an agentic/file-operation test is needed.

If useful, add:

```
run/doctor-ai.sh
```

or:

```
run/doctor-opensource.sh
```

but only if it provides real value.

It should report things like:

```
Ollama       OK  <version>
Qwen Code    OK  <version>
Model        OK  qwen...
Ollama API   OK
VSCodium     OK / not installed
Qwen IDE     OK / manual step required
```

Keep this lightweight.

## 10. VSCodium integration

We use VSCodium, not Microsoft VS Code.

Verify that the CURRENT official Qwen Code Companion supports VS Code forks and is available through Open VSX.

VSCodium uses Open VSX by default.

If `codium` is installed and the Qwen Code Companion extension can reliably be installed from the command line, install it automatically.

Expected shape:

```
codium --install-extension <verified-extension-id>
```

Do NOT guess the extension ID.

Resolve it from the current official Qwen Code extension metadata / Open VSX publication before hardcoding it.

If automatic installation is unreliable, do not hack around it.

Instead print the exact remaining manual steps.

The desired end state is that the user can open Rapidou in VSCodium and interact with Qwen Code natively, including editor context and diffs where supported.

Also investigate the official Qwen commands:

```
/ide install
/ide enable
```

and determine which parts can be automated non-interactively.

If `/ide enable` requires an interactive Qwen session or an IDE-provided environment, leave that as a clearly printed final step instead of trying to fake it.

The installer should end with something like:

```
Local Qwen setup complete.

Model:
  qwen...

CLI:
  qwen

Rapidou/headless:
  <working example command>

VSCodium:
  extension installed

Remaining IDE step:
  1. Restart/open VSCodium in the Rapidou repository
  2. Open a new integrated terminal
  3. Run qwen
  4. Run /ide enable
```

Only print steps that are actually required by the current implementation.

## 11. Do not couple Rapidou to Qwen

Qwen is a first-class agent, not the core architecture.

After these changes it should still be possible to use Rapidou naturally with:

```
Claude Code
Codex
Qwen Code
```

No paid agent should be required.

No Qwen-specific behavior should leak into normal application source code.

Agent integration belongs in the AI/harness/config/scripts layer.

## 12. Keep the implementation small

Avoid architecture like:

```
AgentFactory
AgentProviderManager
LocalAIOrchestrator
BackendRegistry
AdapterStrategyFactory
```

unless the existing code genuinely requires something like that.

Prefer files, conventions, and scripts.

Something like this is conceptually enough:

```
AGENTS.md
CLAUDE.md               # tiny Claude adapter if needed

ai/
  skills/
  hooks/

run/
  install-opensource.sh
  doctor-ai.sh          # only if useful

doc/
  agents.md
```

Use the repository's actual existing layout rather than blindly creating this exact tree.

## 13. Verification

Actually verify the implementation.

At minimum:

* syntax-check shell scripts
* run existing Rapidou tests
* verify no existing Claude/Codex workflow was broken
* verify Qwen configuration format against current official docs
* verify current Ollama model tags
* verify Qwen headless invocation
* verify the Open VSX/VSCodium extension path
* test idempotency where practical
* inspect the final git diff for unnecessary change
