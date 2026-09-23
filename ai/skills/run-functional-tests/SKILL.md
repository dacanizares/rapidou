---
name: run-functional-tests
description: Run and interpret Rapidou's mandatory final functional gate in the Dockerfile Chromium test target. Use after implementation, after review fixes, and before completion; never substitute host-only, API-only, skipped-browser, build-only, or manual checks.
---

# Run functional tests

1. Run exactly `./run/test.sh` from the repository root.
2. Allow Docker, or Podman through the script's Docker-compatible fallback.
3. Require explicit passing results for:
   - `TestAPI`;
   - `TestUI` running real Chromium clicks;
   - `TestSampleMuseum`;
   - `TestDelegationWorkflow`;
   - `TestInstall`;
   - `TestDocumentationRoutes`.
4. Treat `SKIP`, missing Chromium, container failure, or absent `TestUI` output as failure.
5. On failure, report the failing user journey and evidence. Fix only when the current task authorizes implementation, then rerun the complete command.
6. Record the final passing command and results in the handoff.
