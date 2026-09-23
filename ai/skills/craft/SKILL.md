---
name: craft
description: Turn a feature or bug-fix prompt into a completed Rapidou change using the repository's direct spec-as-plan workflow, independently delegated backend/frontend work, mandatory Docker Chromium tests, independent review, and final retest. Use when the user asks to build, change, or fix application behavior end to end.
---

# Craft feature

Treat the spec as the plan. Do not create a second planning artifact.

1. Read `README.md`, `docs/index.md`, and the current relevant spec.
2. Invoke `$spec` to write the happy path, realistic user errors, tests, and acceptance behavior.
3. Invoke `$questions`. Ask one question at a time only while a material gap remains; pause implementation until resolved.
4. Use `$luna` when repository discovery is broad enough to deserve an independent read-only search.
5. Split approved work into explicit backend and frontend ownership:
   - invoke `$select-agent-model` and `./ai/hooks/prepare-agent.sh select-agent-model <codex|claude> ...` before every spawn;
   - delegate `$backend` and `$frontend` to independent agents, in parallel only when their file ownership does not overlap;
   - keep shared spec and coordination files owned by the parent.
6. Reconcile agent output against the spec. Add or correct the functional API and Chromium journeys each domain requires.
7. Invoke `$run-functional-tests`. Do not continue until the complete containerized gate passes.
8. Invoke `$select-agent-model`, then delegate `$code-review` to a fresh independent read-only agent with the diff and spec paths.
9. Fix all valid findings with the smallest direct changes. Update the spec only if observable behavior changed.
10. Invoke `$run-functional-tests` again after every review fix and on the exact final tree.
11. Commit when the user or repository workflow requests it. Report the spec, implementation, review outcome, both functional gates, and commit; never push unless asked.

If backend or frontend work is absent, skip that implementation agent. Never skip spec, questions audit, first test gate, review, or final test gate.
