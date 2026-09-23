---
name: questions
description: Resolve material open questions left by a behavioral spec, one question at a time. Use after the spec skill and before implementation only when an unanswered choice would materially change user-visible behavior, data, security, or scope.
---

# One-by-one questions

1. Read the spec's `Open questions` section.
2. Drop questions answerable from current code, repository rules, or the smallest behavior compatible with the request.
3. Ask the programmer exactly one concise question: the highest-impact unresolved choice.
4. Explain the concrete consequence of the choice in one sentence when useful.
5. Wait for the answer, update the spec, then repeat only if another material question remains.
6. Remove `Open questions` before implementation begins.

Do not bundle questions, ask for preferences that do not change behavior, or turn clarification into a design ceremony.
