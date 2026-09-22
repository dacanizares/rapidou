# Functional harness

The functional harness lives at `tst/`. This Rapidou repository is designed to be consumed by an application repository as a Git submodule, following the same parent/child relationship used by Atlas.

Rapidou owns only `src/functional_test.go`, which creates an isolated application and passes its `http.Handler` to the harness. The harness:

- never imports Rapidou internals;
- never reads Rapidou's SQLite database;
- exercises complete API requests with `httptest`;
- exercises a complete browser flow with `chromedp`;
- travels with Rapidou and is pinned independently by the application's submodule pointer.

Clone restoration will use:

```sh
git submodule update --init --recursive
./run/test
```

The consuming application owns `.gitmodules` and the gitlink that pins Rapidou. This repository does not contain a nested submodule.
