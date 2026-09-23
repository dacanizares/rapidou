---
name: spec
description: Create or update a concise observable-behavior specification that doubles as the implementation plan. Use for new features, behavior changes, or bug fixes before coding; include the happy path, realistic user errors, and testable acceptance behavior without project-management ceremony.
---

# Direct specification

1. Read `README.md`, `docs/index.md`, `docs/shared/structure.md`, `docs/shared/testing.md`, and the relevant current spec.
2. Inspect only enough existing behavior to avoid contradicting reality; use `$luna` for broad discovery.
3. Update the relevant file under `docs/specs/`. The spec is the plan; do not create a separate normal-feature plan.
4. Write concrete statements in this shape: user action → system result.
5. Include:
   - happy path;
   - persistence and permission outcomes when relevant;
   - realistic user mistakes and expected errors;
   - required API and click-driven UI coverage;
   - explicit acceptance conditions.
6. Add a short `Open questions` section only for decisions that materially change behavior and cannot be safely inferred. Remove it when resolved.

Avoid architecture proposals, personas, timelines, estimates, story points, exhaustive edge-case catalogs, and duplicated repository rules.
