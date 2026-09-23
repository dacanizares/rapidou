# Documentation entry points

## Happy path

- A human opens `README.md` and quickly understands what Rapidou is, how to install it at `lib/rapidou`, and what prompt starts a new application.
- Codex or Claude reads `AGENTS.md`/`CLAUDE.md` for mandatory operating rules, then follows `docs/index.md` only to the context required by the task.
- Detailed stack, architecture, frontend, backend, authentication, database, testing, organization, and decision rules remain available under `docs/` without being duplicated in the README.
- The Museo Pixel credentials and local commands remain documented as an example, while `src/` is clearly not application-owned source in a consuming repository.

## Expected mistakes

- Installing the submodule without running `run/install` leaves client discovery incomplete; the README must show both steps.
- Treating Rapidou as a code generator executable is corrected by showing that Codex or Claude runs the `craft` workflow from a prompt.
- Treating the museum as required product behavior is corrected by identifying it as the replaceable base example.

## Acceptance

- README is a short onboarding page rather than the architecture manual.
- No existing project rule or example operating detail is lost; each has a documented destination.
- Documentation links resolve, agent routing tests pass, and `./run/test` passes with the real Chromium journey before and after review.
