# Rapidou functional harness

This directory owns Rapidou's portable functional contract. It sees the application only as an `http.Handler`: API flows use `httptest`, and the UI flow runs the same handler through a real browser with `chromedp`.

Rapidou keeps a thin bridge in `src/functional_test.go`. A consuming application may pin this Rapidou repository as a Git submodule. The harness does not import application internals or inspect its database.

Run it from the parent repository with:

```sh
./run/test
```

The command always builds and runs the Dockerfile `functional-test` target, which includes Chromium. The UI journey is mandatory and fails if the browser cannot start; it must never be skipped as an acceptable final result. Docker is preferred, with Podman supported as a Docker-compatible fallback through the same Dockerfile.
