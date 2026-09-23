# Local agent benchmarks

`run/benchmark-instagram.sh` creates isolated Git worktrees from one Rapidou
commit and gives every trial the same Instagram-like product prompt. It records
each agent's raw machine-readable output and `metadata.json` with the elapsed
wall-clock time and exit status. Model downloads are performed in a separate
preparation step so they are not part of the result.

The three standard trials are deliberately distinct:

| Trial | Agent/tool loop | Model |
|---|---|---|
| `codex` | Authenticated Codex CLI | Your configured Codex default |
| `codex-oss-qwen27b` | Codex CLI using its local Ollama provider | `qwen3.8:27b` |
| `qwen-9b` | Qwen Code CLI | `qwen3.5:9b` |
| `qwen-27b` | Qwen Code CLI | `qwen3.8:27b` |

`codex` is the baseline: it uses no local Qwen model. `codex-oss-qwen27b`
measures the Codex CLI agent loop with an open-weight model. It does not claim
that a cloud Codex model delegates natively to Qwen: that capability does not
exist in the client contract. The two Qwen trials isolate the effect of model
size under Qwen Code's own loop.

## Run

First fetch both models; this can take time and is intentionally excluded from
the benchmark timing:

```sh
./lib/rapidou/run/benchmark-instagram.sh --prepare
```

From `run/`, launch everything with one command:

```sh
./benchmark-instagram.sh
```

It creates `run/tmp/only-codex`, `run/tmp/codex-ollama`, `run/tmp/qwen-9b`, and
`run/tmp/qwen-27b`. Each contains `app/` (the resulting application worktree),
the immutable `prompt.md`, `agent.json` or `agent.jsonl`, and `metadata.json`.
The script refuses to overwrite an existing `run/tmp` benchmark. The defaults
allow 80 Qwen turns and 90 minutes per trial; adjust them only for all trials if
you need a different budget:

```sh
./benchmark-instagram.sh --turns 80 --wall-time 90m
```

Run one trial only with `--run codex`, `--run qwen-9b`, `--run qwen-27b`, or
`--run codex-oss-qwen27b`. Review the four `metadata.json` files first, then
open each `app/` and run its documented functional test before judging the
generated application.
