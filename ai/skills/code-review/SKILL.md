---
name: code-review
description: Review a Rapidou change against its behavioral spec, functional tests, and simplicity rules. Use after the first complete test pass and before final verification, preferably through an independent read-only agent.
---

# Rapidou code review

1. Read `README.md`, `docs/shared/index.md`, the relevant domain index, the spec, and the complete diff.
2. Review for observable bugs first:
   - behavior diverging from spec;
   - missing happy path or realistic user-error coverage;
   - auth, validation, persistence, migration, or state regressions;
   - UI click, keyboard, dialog, responsive, or stale-render problems.
3. Review simplicity second:
   - unnecessary layers, dependencies, indirection, files, or ceremonies;
   - duplicated behavior or documentation;
   - hidden control flow that direct functions would clarify.
4. Return findings ordered by severity with exact paths and lines. State `No findings` when appropriate.
5. Do not edit during the independent review. The parent owns fixes and decides whether a finding changes the spec.
