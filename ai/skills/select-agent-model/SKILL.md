---
name: select-agent-model
description: Select the Codex or Claude model for every independent subagent from the task's complexity and size. Use immediately before any subagent spawn, whether sequential or parallel; delegation is forbidden until this skill and the repository delegation hook have been used.
---

# Select agent model

1. Summarize the delegated task in one sentence with explicit ownership.
2. Classify complexity:
   - `low`: lookup, inventory, narrow mechanical change.
   - `medium`: normal feature, test, or review requiring local judgment.
   - `high`: cross-cutting design, security, concurrency, subtle debugging, or broad review.
3. Classify size:
   - `small`: one focused result or roughly 1–2 files.
   - `medium`: several connected results or roughly 3–8 files.
   - `large`: broad repository surface or more than 8 files.
4. Identify the active client and run `./ai/hooks/prepare-agent.sh select-agent-model <codex|claude> <complexity> <size>`.
5. Use its exact `model` in the spawn call. For Codex, also use the exact `reasoning_effort`. Pass only the context the agent needs and name its writable scope.
6. If the task changes materially, classify again before sending follow-up work.

Never select a model from agent availability, convenience, or parallelism alone. Complexity and size decide.
