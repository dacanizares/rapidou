# Backend context

Read [shared context](../shared/index.md) first.

Backend code lives in `src/*.go`; SQL assets live in `data/sql/`; API functional coverage lives in `tst/api.go` and is connected through `src/functional_test.go`.

Use `net/http`, `database/sql`, SQLite, direct handlers, direct database functions, and explicit validation. Test complete HTTP behavior with `httptest` and a real temporary SQLite database. Cover the happy path plus realistic authentication, validation, missing-resource, conflict, and persistence failures relevant to the spec.
