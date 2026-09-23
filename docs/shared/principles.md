# Application principles

Rapidou favors the simplest complete design that works, remains versatile, and can be understood from top to bottom by humans and machines.

## Deliberate stack

Use:

- Go and its standard library whenever practical;
- `net/http` and `database/sql`;
- SQLite with direct SQL;
- JWT authentication and bcrypt passwords;
- embedded HTML, CSS, and plain JavaScript;
- `testing`, `httptest`, and `chromedp`;
- Docker with Ubuntu support;
- one runnable application binary.

Do not add Node.js, npm, a frontend build system, or a frontend framework. Add dependencies only when a concrete feature cannot reasonably be implemented with the language or browser primitives.

## Direct architecture

Build primarily from structs, functions, direct calls, plain data, simple control flow, and simple composition. Database functions may be grouped for navigation; handlers call them directly.

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

When two behaviors exist, prefer explicit control flow:

```go
if condition {
    result = functionA()
} else {
    result = functionB()
}
```

Do not introduce services, repositories, use cases, adapters, ports, factories, controller layers, dependency-injection frameworks, domain layers, strategies, registries, providers, or generic infrastructure merely to represent direct behavior.

Functions are the main unit of behavior. Structs are the main unit of data. A clear flow is:

```text
HTTP request
→ handler function
→ validation function
→ database function
→ response function
```

Do not hide simple behavior behind an abstraction.

## Backend

Use `http.ServeMux` and normal handlers:

```go
mux.HandleFunc("GET /api/users", auth(listUsers))
mux.HandleFunc("POST /api/users", auth(createUser))
```

Small concrete helpers such as `readJSON`, `writeJSON`, `writeError`, and `authenticate` are useful. Do not build an internal framework.

Keep authentication in a few functions. Browser authentication should prefer an HttpOnly, Secure, SameSite cookie; the API may also accept `Authorization: Bearer <token>`.

Use SQLite through `database/sql`. Write SQL directly and prefer functions such as `userByID`, `userByEmail`, `insertUser`, `updateUser`, `deleteUser`, and `listUsers`. Do not add an ORM unless a current requirement proves direct SQL insufficient.

## Frontend

Use browser APIs directly. State can be a plain object, API calls plain async functions, and rendering plain functions:

```js
const state = {
    users: [],
    currentUser: null
};

async function refreshUsers() {
    state.users = await api("/api/users");
    renderUsers("#users");
}
```

Render only what changed when practical. Prefer functions such as `renderUsers`, `renderSettings`, `renderNavigation`, `showDialog`, and `hideDialog` over a rendering framework or component system.

UI elements may simply be functions that return DOM elements.

Use `fetch`, `querySelector`, `createElement`, `replaceChildren`, `addEventListener`, and browser history APIs when routing is needed. Do not add a virtual DOM, JSX, a reactive framework, or a state-management library.

## Functional testing

Tests follow the same direct style. API tests use `httptest` with helpers such as `get`, `postJSON`, `login`, `expectStatus`, and `expectJSON`.

UI tests use `chromedp` with small helpers such as `open`, `click`, `typeText`, `waitVisible`, `waitText`, `expectText`, and `eventually`. They read as user actions:

```go
func TestCreateUserUI(t *testing.T) {
    open(ctx, "/users")
    click(ctx, "#new-user")
    typeText(ctx, "#name", "Daniel")
    click(ctx, "#save")
    expectText(ctx, "#users", "Daniel")
}
```

Avoid testing frameworks, DSLs, page-object architectures, BDD systems, and test abstraction layers unless they become genuinely necessary. The complete Docker/Chromium gate is defined in [testing](testing.md) and [the harness](../harness/index.md).

## Organization and decisions

Start small. Keep Go application files in `package main`; split files only when navigation improves. A small application can remain close to:

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

Choose the simplest design that works correctly. Prefer explicit code over implicit behavior, direct calls over indirection, concrete code over generic infrastructure, browser and language primitives over frameworks, and a small amount of boring code over a sophisticated abstraction.

Do not solve hypothetical future problems or add infrastructure for possible future requirements. Implement the feature that exists now. The goal is not minimal code at any cost; it is simple, complete, adaptable code.

## Code clarity

Write code that is easy for humans and AI agents to inspect, modify, and debug.

Prefer:

- boring > clever
- explicit > compressed
- local reasoning > abstraction
- simple control flow > dense expression chains
- named intermediate values > deeply nested expressions

Use loops and `if` statements when they make multi-step logic clearer.

Avoid clever one-liners, long method chains, unnecessary abstractions, metaprogramming, and hidden control flow.

Do not make simple idiomatic code artificially verbose. Simple idioms are fine when they represent one obvious operation:

```ts
const user = users.find((user) => user.id === userId);
```

For multi-step logic, prefer explicit code:

```ts
// BAD
const result = Array.from(new Set(array)).map((x) => x * 2).filter((x) => x > 10);

// GOOD
const result: number[] = [];
const seen = new Set<number>();

for (const x of array) {
    if (seen.has(x)) {
        continue;
    }

    seen.add(x);

    const doubled = x * 2;

    if (doubled > 10) {
        result.push(doubled);
    }
}
```
