# rapidou

Rapidou is a deliberately small Go application base. The included example is **Museo Pixel**, a public video game collection with authenticated curator tools, SQLite persistence, JWT/bcrypt authentication, an embedded plain-JavaScript frontend, and a portable functional harness.

Run it locally:

```sh
./run/dev
```

Then open `http://localhost:8080` and sign in with the development account:

```text
admin@rapidou.test
rapidou-local-password
```

Run the complete functional contract and build the single binary with:

```sh
./run/test
./run/build
```

The browser flow runs when Chrome or Chromium is installed. API flows always run with `httptest`. Production startup requires an explicit `APP_JWT_SECRET`; the defaults above exist only in `run/dev`.

Build this application with extreme simplicity as the primary architectural constraint.

The application must remain plain, direct, small, and easy for both humans and machines to understand.

Use:

* Go.
* Go standard library whenever practical.
* `net/http` for HTTP and routing.
* `database/sql`.
* SQLite.
* JWT authentication.
* bcrypt for passwords.
* HTML, CSS, and plain JavaScript.
* No Node.js.
* No npm.
* No frontend build system.
* No frontend framework.
* `testing` and `httptest` for API tests.
* `chromedp` for browser/UI functional tests.
* Docker with Ubuntu support.
* Embed the frontend into the Go executable.
* Produce a single runnable application binary.

## Core architecture

The application should be built almost entirely from:

* structs,
* functions,
* direct function calls,
* plain data,
* simple control flow,
* simple composition of functions.

Avoid architectural patterns unless a concrete problem requires them.

Do not introduce layers such as:

* services,
* repositories,
* use cases,
* adapters,
* ports,
* factories,
* controllers as a separate abstraction,
* dependency injection frameworks,
* domain layers,
* generic abstraction layers.

Database access functions may be grouped together for organization.

Endpoint handlers should call those functions directly.

Example:

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    input := readInput(r)

    user, err := createUser(db, input)

    if err != nil {
        writeError(w, err)
        return
    }

    writeJSON(w, http.StatusCreated, user)
}
```

If there are two ways to perform something, prefer simple control flow:

```go
if condition {
    result = functionA()
} else {
    result = functionB()
}
```

Do not create an interface, strategy pattern, factory, registry, provider, or abstraction merely to represent this choice.

## Functions first

Functions are the main unit of behavior.

Structs are the main unit of data.

Prefer plain function composition:

```text
HTTP request
→ handler function
→ validation function
→ database function
→ response function
```

A function calling another function directly is good architecture.

Do not hide simple behavior behind unnecessary abstractions.

## Frontend

Use plain browser APIs.

State can be a plain object:

```js
const state = {
    users: [],
    currentUser: null
};
```

API calls are plain async functions.

Rendering is done by functions.

Render only what changed when practical.

Example:

```js
async function refreshUsers() {
    state.users = await api("/api/users");
    renderUsers("#users");
}
```

Prefer:

```js
renderUsers()
renderSettings()
renderNavigation()
showDialog()
hideDialog()
```

over introducing a rendering framework or component system.

UI elements may simply be functions returning DOM elements.

Use:

* `fetch`
* `querySelector`
* `createElement`
* `replaceChildren`
* `addEventListener`
* browser history APIs when routing is needed.

No virtual DOM.

No JSX.

No reactive framework.

No state-management library.

No frontend dependency unless there is a concrete feature that cannot reasonably be implemented with browser APIs.

## Backend

Prefer the Go standard library.

Use `http.ServeMux`.

Use normal handlers:

```go
mux.HandleFunc("GET /api/users", auth(listUsers))
mux.HandleFunc("POST /api/users", auth(createUser))
```

Use small helpers where useful:

```go
readJSON()
writeJSON()
writeError()
authenticate()
```

Keep helpers obvious and concrete.

Do not build an internal framework.

## Authentication

Use JWT.

Prefer an HttpOnly, Secure, SameSite cookie for browser authentication.

The API may additionally support:

```text
Authorization: Bearer <token>
```

Keep authentication logic centralized in a few simple functions.

Do not create an authentication framework.

## Database

Use SQLite through `database/sql`.

Write SQL directly.

Prefer functions such as:

```go
userByID()
userByEmail()
insertUser()
updateUser()
deleteUser()
listUsers()
```

Do not use an ORM unless a concrete requirement later proves that raw SQL is insufficient.

## Testing

Tests must also follow the same philosophy.

API functional tests should be simple Go functions using `httptest`.

Create small helpers such as:

```go
get()
postJSON()
login()
expectStatus()
expectJSON()
```

UI tests should use `chromedp`.

Build simple browser helpers such as:

```go
open()
click()
typeText()
waitVisible()
waitText()
expectText()
eventually()
```

Tests should read like sequences of actions.

Example:

```go
func TestCreateUserUI(t *testing.T) {
    open(ctx, "/users")
    click(ctx, "#new-user")
    typeText(ctx, "#name", "Daniel")
    click(ctx, "#save")
    expectText(ctx, "#users", "Daniel")
}
```

Avoid testing frameworks, DSLs, page-object architectures, BDD systems, or test abstraction layers unless they become genuinely necessary.

## Project organization

Start small.

Prefer something close to:

```text
main.go
auth.go
db.go
users.go
main_test.go
ui_test.go

web/
    index.html
    app.js
    app.css

schema.sql
go.mod
Dockerfile
```

All Go application files may remain in `package main`.

Do not split the application into packages merely for architectural appearance.

Split files only when it improves navigation.

## Decision rule

When implementing something, choose the simplest design that works correctly.

The code should be obvious enough that a machine can understand the flow without reconstructing hidden architectural conventions.

Prefer explicit code over implicit behavior.

Prefer direct calls over indirection.

Prefer concrete code over generic infrastructure.

Prefer browser and language primitives over frameworks.

Prefer a small amount of boring code over a sophisticated abstraction.

Do not solve hypothetical future problems.

Do not add infrastructure for possible future requirements.

Implement the feature that exists now.

The target is not minimal code at any cost.

The target is:

**the simplest complete code that works, remains versatile, and can be understood from top to bottom.**
