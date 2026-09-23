# Rapidou functional harness

This directory owns Rapidou's portable functional contract. It sees the application only as an `http.Handler`: API flows use `httptest`, and the UI flow runs the same handler through a real browser with `chromedp`.

The base example keeps a thin bridge in `src/functional_test.go`. A consuming application pins this repository at `lib/rapidou` and adapts that bridge in its own source. The harness does not import application internals or inspect its database.

Run it from the parent repository with:

```sh
./lib/rapidou/run/install
./run/test
```

The command always builds and runs the Dockerfile `functional-test` target, which includes Chromium. The UI journey is mandatory and fails if the browser cannot start; it must never be skipped as an acceptable final result. Docker is preferred, with Podman supported as a Docker-compatible fallback through the same Dockerfile.
