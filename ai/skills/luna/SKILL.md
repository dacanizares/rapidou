---
name: luna
description: Perform focused read-only repository discovery through a lower-complexity independent subagent. Use to locate files, call paths, tests, specs, or likely change surfaces before implementation, especially when the main agent would otherwise spend context on broad code search.
---

# Luna code search

1. Define one bounded search question and the evidence required.
2. Invoke `$select-agent-model` and the delegation hook. A normal Luna task is `low small`; expand only when the search surface truly requires it.
3. Spawn one independent read-only agent with minimal context and `fork_turns="none"` when available.
4. Require the agent to use `rg`/`rg --files`, make no edits, and return:
   - relevant paths and line numbers;
   - the observed call/data flow;
   - unanswered facts, clearly labeled.
5. Verify important findings in the main thread before editing.

Do not use Luna for product decisions, implementation, or speculative redesign.
