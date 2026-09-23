---
name: backend
description: Implement the backend portion of an approved Rapidou specification using direct Go, net/http, SQLite, and functional API tests. Use for routes, authentication, validation, persistence, schema, seed, or server behavior changes.
---

# Rapidou backend

1. Read `docs/shared/index.md`, `docs/backend/index.md`, and the relevant spec.
2. Inspect adjacent handlers, database functions, routes, and tests before editing.
3. Implement only the specified behavior with plain structs, functions, direct calls, explicit validation, and plain SQL.
4. Keep backend ownership to `src/*.go`, `data/sql/`, and API coverage in `tst/api.go` unless the parent assigns another boundary.
5. Add functional tests through the real HTTP handler and temporary SQLite database:
   - every specified happy path;
   - realistic user input, authentication, conflict, missing-resource, and persistence errors that apply.
6. Preserve existing databases with the smallest explicit migration needed.
7. Report changed behavior, files, and any frontend contract the parent must consume.

Do not add services, repositories, ORMs, dependency injection, generic frameworks, or unit tests for simple internal functions.
