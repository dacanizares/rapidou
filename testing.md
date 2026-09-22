## Testing philosophy

Testing should prioritize real application behavior, not implementation details.

The goal is functional confidence, not line coverage.

Do not optimize tests for code coverage percentages.

Do not create large numbers of tests merely to exercise individual branches or lines.

Test useful behavior.

Think primarily in terms of:

* what a real user wants to accomplish,
* what mistakes a real user may make,
* what invalid states can occur,
* what the application must guarantee,
* what failures would actually matter.

A smaller number of meaningful end-to-end functional tests is preferred over hundreds of narrow tests.

## API functional tests

API tests should exercise complete endpoints.

A test should send a real HTTP request through the application's HTTP handler and verify the complete observable result.

Prefer tests such as:

```text
create user successfully
reject invalid user
reject duplicate email
login successfully
reject wrong password
reject unauthenticated request
update existing resource
reject update of missing resource
delete resource
reject forbidden action
```

The test should exercise the same path used by the real application:

```text
HTTP request
→ routing
→ authentication
→ validation
→ handler
→ database
→ HTTP response
```

Do not bypass parts of the application merely to make tests easier.

Use `httptest` and a real temporary SQLite database when practical.

Tests should verify meaningful results:

* HTTP status,
* response body,
* persisted state,
* authentication behavior,
* important side effects.

Example:

```go
func TestCreateUser(t *testing.T) {
    app := newTestApp(t)

    token := login(t, app, "admin@test.com", "secret")

    response := postJSON(t, app, "/api/users", token, map[string]any{
        "name":  "Daniel",
        "email": "daniel@test.com",
    })

    expectStatus(t, response, http.StatusCreated)
    expectJSONField(t, response, "email", "daniel@test.com")

    user := getUserByEmail(t, app.DB, "daniel@test.com")
    assertEqual(t, user.Name, "Daniel")
}
```

Also test negative behavior:

```go
func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
    ...
}
```

Negative cases are important because they represent mistakes and invalid actions that real users can make.

Do not attempt to enumerate every theoretically possible invalid input.

Test the failures that are realistic, important, or dangerous.

## UI functional tests

UI tests should exercise complete user flows in a real browser.

The full UI flow is a mandatory completion gate. Run it through `./run/test`, inside the Docker test target that owns Chromium. Never accept a skipped browser test or an API-only run as final verification.

Think in terms of user journeys rather than individual buttons or components.

Prefer flows such as:

```text
login
→ open users
→ create user
→ see new user
```

or:

```text
login
→ enter invalid data
→ submit
→ see useful validation error
→ correct data
→ submit successfully
```

or:

```text
open protected page while logged out
→ redirected to login
→ authenticate
→ continue successfully
```

UI tests should include both positive and negative flows.

Positive flows verify that users can accomplish important tasks.

Negative flows verify that the application behaves correctly when users:

* enter invalid values,
* omit required values,
* use incorrect credentials,
* repeat an action,
* attempt something unauthorized,
* navigate unexpectedly,
* interact with stale or missing data.

Do not create a UI test for every element.

Do not test CSS implementation details.

Do not test individual DOM functions in isolation.

Test what the user can actually observe and do.

UI tests should read like simple sequences of actions:

```go
func TestCreateUserFlow(t *testing.T) {
    open(ctx, "/login")

    typeText(ctx, "#email", "admin@test.com")
    typeText(ctx, "#password", "secret")
    click(ctx, "#login")

    waitURL(ctx, "/users")

    click(ctx, "#new-user")
    typeText(ctx, "#name", "Daniel")
    typeText(ctx, "#email", "daniel@test.com")
    click(ctx, "#save")

    expectText(ctx, "#users", "Daniel")
}
```

Negative example:

```go
func TestCreateUserShowsValidationError(t *testing.T) {
    loginUI(ctx)

    open(ctx, "/users")
    click(ctx, "#new-user")
    click(ctx, "#save")

    expectText(ctx, "#user-form", "Email is required")
}
```

Build only a small set of simple browser helpers:

```go
open()
click()
typeText()
waitVisible()
waitHidden()
waitText()
waitURL()
expectText()
eventually()
```

Tests should remain plain functions composed from these helpers.

## No unit tests by default

Do not create unit tests by default.

Do not test internal functions merely because they exist.

Internal implementation is allowed to change without forcing unrelated tests to change.

Functional tests are the primary testing strategy.

Add a unit test only if a specific piece of isolated logic becomes sufficiently complex that testing it through the public behavior would be impractical.

This should be exceptional, not the default.

## Assertions and invariants

Prefer explicit runtime assertions for internal programming invariants when useful.

Assertions belong at boundaries where violating a condition indicates a programming bug.

Examples:

```go
func saveUser(db *sql.DB, user User) error {
    assert(db != nil)
    assert(user.ID != 0)

    ...
}
```

or:

```go
func renderUsers(selector string) {
    console.assert(selector !== "")
    ...
}
```

Assertions should express assumptions that our own code must guarantee.

Do not use assertions for normal user errors or expected runtime failures.

For example:

```text
empty email
wrong password
missing record
invalid request
expired token
permission denied
```

must be handled normally and returned as proper application errors.

Assertions are for situations such as:

```text
this database connection must exist
this internal ID must already have been assigned
this state should be impossible here
this function must never be called without initialization
```

If such an invariant fails, failing loudly is preferable to continuing in an invalid internal state.

## Test selection rule

Before adding a test, ask:

> What real failure or user behavior does this test protect us from?

If there is no meaningful answer, the test probably does not need to exist.

The testing target is:

**functional coverage of important behavior, not numerical coverage of implementation.**
