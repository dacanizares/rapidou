# Install as a submodule

From the consuming repository:

```sh
git submodule add <rapidou-repository-url> lib/rapidou
./lib/rapidou/run/install
```

Commit `.gitmodules`, the `lib/rapidou` gitlink, `.agents/`, `.claude/`, `.codex/`, `AGENTS.md`, and `CLAUDE.md`.

The installer links one canonical set of skills into the paths Codex and Claude discover, installs the agent-spawn guard, and adds short shared instructions. It is idempotent and refuses to overwrite unrelated files. If `.codex/hooks.json` or `.claude/settings.json` already exists, merge the matching fragment from `lib/rapidou/ai/install/`, inspect it, then rerun with `--accept-existing-hooks`.

Restart both clients after installation. In Codex, open `/hooks` and trust the project hook; Codex intentionally does not run an unreviewed repository hook.

`lib/rapidou/src/` is a working base example: copy ideas or start the application from it deliberately. The installer never copies it. The consuming application must adapt `lib/rapidou/src/functional_test.go` to its own handler and own a `./run/test` that runs the complete harness in Docker with Chromium.

After cloning:

```sh
git submodule update --init --recursive
./lib/rapidou/run/install
./run/test
```
