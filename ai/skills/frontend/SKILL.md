---
name: frontend
description: Implement the frontend portion of an approved Rapidou specification with plain HTML, CSS, JavaScript, and click-driven Chromium tests. Use for rendering, forms, dialogs, interactions, responsive behavior, accessibility, or other user-visible changes.
---

# Rapidou frontend

1. Read `docs/shared/index.md`, `docs/frontend/index.md`, and the relevant spec.
2. Inspect the current DOM, state, API calls, rendering functions, styles, and browser journey.
3. Implement with browser primitives in `src/web/`; keep content in HTML/JS and normal presentation in CSS.
4. Preserve keyboard access, responsive layout, clear focus, and independent behavior for nested controls.
5. Extend `tst/ui.go` with a real user journey that covers:
   - the specified happy path with clicks and visible outcomes;
   - realistic user mistakes from the spec;
   - regression assertions for dialog, navigation, and persisted state behavior that changed.
6. Report changed behavior, files, and any backend contract assumption.

Do not add Node.js, npm, frameworks, build tools, CSS utility systems, inline-style generation, or isolated DOM unit tests.
