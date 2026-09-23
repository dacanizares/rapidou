# Functional harness

The functional harness lives at `tst/`. This repository can be pinned at `lib/rapidou`; see the [installation guide](../installation.md).

The base example owns `src/functional_test.go`, which creates an isolated application and passes its `http.Handler` to the harness. A consuming application adapts that thin bridge in its own source. The harness:

- never imports Rapidou internals;
- never reads Rapidou's SQLite database;
- exercises complete API requests with `httptest`;
- verifies authentication, user administration, and the public museum contract;
- exercises login and museum cataloging in a real browser with `chromedp`;
- builds a dedicated Docker test target that includes Chromium;
- fails instead of skipping when a real browser cannot be started;
- travels with Rapidou and is pinned independently by the application's submodule pointer.

Restore a consuming clone with:

```sh
git submodule update --init --recursive
./lib/rapidou/run/install
./run/test
```

`./run/test` is the mandatory completion gate for every implementation change. It selects Docker first and Podman only as a Docker-compatible fallback, builds the `functional-test` target, and runs the complete API and browser journeys inside that container. A host-only `go test`, an API-only result, or a skipped UI test does not satisfy the gate.

The consuming application owns `.gitmodules` and the gitlink that pins Rapidou. This repository does not contain a nested submodule.

## Agent workflow

Reusable workflow skills live under `ai/skills/` and are indexed at `ai/skills/index.md`. The installer links them into both clients' discovery directories. Before every independent agent spawn, read `select-agent-model` and run:

```sh
./ai/hooks/prepare-agent.sh select-agent-model <codex|claude> <complexity> <size>
```

The returned model and, for Codex, reasoning effort are mandatory for that spawn. A one-shot client hook enforces the selection. This applies to parallel work, sequential work, Luna searches, and code review. `craft` owns the complete spec → questions → backend/frontend → test → review → retest sequence.
