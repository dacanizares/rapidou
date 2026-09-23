# Frontend context

Read [shared context](../shared/index.md) and [styling](styling.md) first.

Frontend code lives in `src/web/index.html`, `src/web/app.js`, and `src/web/app.css`. Use plain HTML, CSS, JavaScript, and native browser APIs. Do not add Node.js, npm, a frontend framework, or a build step.

User-visible behavior must be covered by the click-driven Chromium journey in `tst/ui.go`, including the happy path and realistic user mistakes from the spec.
