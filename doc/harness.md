# Functional harness

The functional harness lives at `tst/`. This Rapidou repository is designed to be consumed by an application repository as a Git submodule, following the same parent/child relationship used by Atlas.

Rapidou owns only `src/functional_test.go`, which creates an isolated application and passes its `http.Handler` to the harness. The harness:

- never imports Rapidou internals;
- never reads Rapidou's SQLite database;
- exercises complete API requests with `httptest`;
- verifies authentication, user administration, and the public museum contract;
- exercises login and museum cataloging in a real browser with `chromedp`;
- builds a dedicated Docker test target that includes Chromium;
- fails instead of skipping when a real browser cannot be started;
- travels with Rapidou and is pinned independently by the application's submodule pointer.

Clone restoration will use:

```sh
git submodule update --init --recursive
./run/test
```

`./run/test` is the mandatory completion gate for every implementation change. It selects Docker first and Podman only as a Docker-compatible fallback, builds the `functional-test` target, and runs the complete API and browser journeys inside that container. A host-only `go test`, an API-only result, or a skipped UI test does not satisfy the gate.

The consuming application owns `.gitmodules` and the gitlink that pins Rapidou. This repository does not contain a nested submodule.
