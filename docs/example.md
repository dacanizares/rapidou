# Base example: Museo Pixel

`src/` is a complete example application, not code that the installer copies into a consuming repository. It demonstrates a public video-game museum with detail views, Steam/GOG links, popup authentication, responsive add/edit dialogs, URL and uploaded images, SQLite persistence, JWT/bcrypt authentication, an embedded plain-JavaScript frontend, and the portable functional harness.

Run it from this repository:

```sh
./run/dev
```

Open `http://localhost:8080` and use:

```text
admin@rapidou.test
rapidou-local-password
```

Build and verify it with:

```sh
./run/build
./run/test
```

Production startup requires an explicit `APP_JWT_SECRET`; the defaults above exist only in `run/dev`.

The consuming application owns its adapted source and its `src/functional_test.go` bridge. Use this example to see the intended scale and direct structure, then implement the requested application rather than preserving museum-specific behavior.
