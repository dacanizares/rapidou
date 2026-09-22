# Rapidou functional harness

This repository owns Rapidou's portable functional contract. It sees an application only as an `http.Handler`: API flows use `httptest`, and the UI flow runs the same handler through a real browser with `chromedp`.

Rapidou keeps a thin bridge in `src/functional_test.go`. A consuming application may pin this Rapidou repository as a Git submodule. The harness must not import application source or inspect its database.

Run it from the parent repository with:

```sh
./run/test
```

Chrome/Chromium is optional for local API-only work. When it is unavailable, the UI flow reports a skip rather than hiding an API failure.
