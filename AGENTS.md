# Project context

Rapidou is a small foundation and functional harness for applications built from Go, SQLite, JWT/bcrypt, embedded plain JavaScript, and direct functions. `src/` is the executable Museo Pixel base example; in a consuming repository, application code belongs to the parent and Rapidou lives at `lib/rapidou`.

Before changing code, read `docs/index.md`, `docs/shared/principles.md`, and the relevant specification. Keep `README.md` limited to human onboarding; put durable operating rules here and detailed guidance under `docs/`.

# Completion gate

Before treating any implementation change as complete, always run:

```sh
./run/test.sh
```

This is mandatory. The command must execute the complete functional harness inside the Dockerfile `functional-test` target (using Docker, or Podman as a Docker-compatible fallback), including the real Chromium UI journey with clicks. An API-only run, a skipped browser test, a successful build, or a manual HTTP check is not sufficient for completion.

Do not bypass, skip, or replace the browser journey because Chrome/Chromium is unavailable on the host. The test container owns Chromium specifically so the UI flow always runs.

# Delegation gate

Before spawning any independent agent, sequentially or in parallel:

1. Read and invoke `ai/skills/select-agent-model/SKILL.md`.
2. Classify the subtask by complexity and size.
3. Identify the active client and run `./ai/hooks/prepare-agent.sh select-agent-model <codex|claude> <complexity> <size>`.
4. Use the exact model returned by the hook and, for Codex, its reasoning effort.
5. Give the agent a bounded task, required skill, context paths, and file ownership.

The client `PreToolUse` hook blocks an agent spawn without a fresh matching selection. There are no exceptions for read-only searches or reviews. Use `ai/skills/luna/SKILL.md` for delegated code search. Use `ai/skills/craft/SKILL.md` for end-to-end feature work.
