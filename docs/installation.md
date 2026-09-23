# Install as a submodule

From the consuming repository:

```sh
git submodule add <rapidou-repository-url> lib/rapidou
./lib/rapidou/run/install.sh
```

Commit `.gitmodules`, the `lib/rapidou` gitlink, `.agents/`, `.claude/`, `.codex/`, `.qwen/`, `AGENTS.md`, and `CLAUDE.md`.

The installer links one canonical set of skills into the paths Codex, Claude Code, and Qwen Code discover, installs each native agent-spawn guard, and adds short shared instructions. It is idempotent and refuses to overwrite unrelated files. If `.codex/hooks.json`, `.claude/settings.json`, or `.qwen/settings.json` already exists, merge the matching fragment from `lib/rapidou/ai/install/`, inspect it, then rerun with `--accept-existing-hooks`.

Restart the clients after installation. In Codex, open `/hooks` and trust the project hook; Codex intentionally does not run an unreviewed repository hook.

For a subscription-free local Qwen Code and Ollama setup, run `./lib/rapidou/run/install-opensource.sh`. See [agent compatibility](agents.md) for model selection, diagnostics, headless use, and VSCodium integration.

`lib/rapidou/src/` is a working base example: copy ideas or start the application from it deliberately. The installer never copies it. The consuming application must adapt `lib/rapidou/src/functional_test.go` to its own handler and own a `./run/test.sh` that runs the complete harness in Docker with Chromium.

After cloning:

```sh
git submodule update --init --recursive
./lib/rapidou/run/install.sh
./run/test.sh
```
