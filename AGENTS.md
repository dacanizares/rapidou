# Completion gate

Before treating any implementation change as complete, always run:

```sh
./run/test
```

This is mandatory. The command must execute the complete functional harness inside the Dockerfile `functional-test` target (using Docker, or Podman as a Docker-compatible fallback), including the real Chromium UI journey with clicks. An API-only run, a skipped browser test, a successful build, or a manual HTTP check is not sufficient for completion.

Do not bypass, skip, or replace the browser journey because Chrome/Chromium is unavailable on the host. The test container owns Chromium specifically so the UI flow always runs.
