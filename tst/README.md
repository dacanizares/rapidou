# Rapidou functional harness

This directory owns Rapidou's portable functional contract. It sees the application only as an `http.Handler`: API flows use `httptest`, and the UI flow runs the same handler through a real browser with `chromedp`.

- `shared.go` contains the common HTTP test types and helpers.
- `api.go` and `ui.go` exercise application behavior.
- `installation_test.go` verifies Rapidou itself: skill discovery, model-selection hooks, installation, and documentation routes.

Run only those installation checks with:

```sh
./run/installation_test.sh
```

This is a focused developer command, not the completion gate.

The base example keeps a thin bridge in `src/functional_test.go`. A consuming application pins this repository at `lib/rapidou` and adapts that bridge in its own source. The harness does not import application internals or inspect its database.

Run it from the parent repository with:

```sh
./lib/rapidou/run/install.sh
./run/test.sh
```

The command always builds and runs the Dockerfile `functional-test` target, which includes Chromium. The UI journey is mandatory and fails if the browser cannot start; it must never be skipped as an acceptable final result. Docker is preferred, with Podman supported as a Docker-compatible fallback through the same Dockerfile.
