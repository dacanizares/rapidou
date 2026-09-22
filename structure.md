## Project structure

Keep the project direct, flat, and practical.

The structure should make it obvious where to find:

* application source,
* tests,
* executable scripts,
* database assets,
* specifications and plans,
* AI instructions,
* runtime data,
* build outputs.

Use:

```text
project/
│
├── src/
│   ├── main.go
│   ├── auth.go
│   ├── db.go
│   ├── users.go
│   ├── ...
│   │
│   └── web/
│       ├── index.html
│       ├── app.js
│       └── app.css
│
├── tst/
│   ├── api_test.go
│   ├── ui_test.go
│   └── helpers.go
│
├── run/
│   ├── dev.sh
│   ├── build.sh
│   ├── test.sh
│   └── docker.sh
│
├── data/
│   ├── app.db
│   │
│   └── sql/
│       ├── schema.sql
│       ├── seed.sql
│       └── scripts/
│           ├── reset.sql
│           ├── sample-data.sql
│           └── ...
│
├── doc/
│   ├── spec.md
│   ├── plan.md
│   └── ...
│
├── ai/
│   ├── skills/
│   │   └── ...
│   │
│   └── hooks/
│       └── ...
│
├── bin/
│   └── app
│
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

---

# `src/`

Contains application source code.

Keep the application as one Go application and normally one Go package:

```go
package main
```

Do not introduce packages only to represent architectural concepts.

Files exist only to make navigation easier.

Example:

```text
main.go
    startup
    configuration
    route registration
    application initialization

auth.go
    login
    logout
    JWT
    password handling
    auth middleware/functions

db.go
    database initialization
    shared database helpers

users.go
    user endpoints
    user validation
    user database operations
```

Functions may call functions from other files directly.

Example:

```text
createUserHandler()
    ↓
validateUser()
    ↓
insertUser()
    ↓
writeJSON()
```

This is good architecture.

There is no need for:

```text
controller
service
repository
use case
adapter
provider
factory
domain layer
```

when direct function composition already solves the problem.

If there are two behaviors:

```go
if condition {
    result = functionA()
} else {
    result = functionB()
}
```

Prefer that over introducing an abstraction representing the choice.

---

# `src/web/`

Contains the complete frontend.

Use:

```text
index.html
app.js
app.css
```

Add files only when the existing files become genuinely difficult to navigate.

Do not introduce:

```text
Node.js
npm
React
Vue
Angular
Tailwind
frontend bundlers
CSS preprocessors
frontend compilation
```

by default.

Frontend behavior should be plain JavaScript functions.

Example:

```js
async function refreshUsers() {
    state.users = await api("/api/users");
    renderUsers("#users");
}
```

Rendering should normally be local.

If users changed:

```js
renderUsers("#users");
```

If navigation changed:

```js
renderNavigation("#navigation");
```

If a dialog changed:

```js
renderUserDialog("#user-dialog");
```

Do not create a global render cycle unless the application actually needs one.

CSS should use:

```text
CSS variables
classes
Flexbox
Grid
media queries
native browser features
```

Prefer:

```css
.btn
.btn-primary
.card
.grid
.stack

.section-users
.users-toolbar
.user-row
```

over inline styles or large utility-class vocabularies.

JavaScript should normally manipulate:

```text
content
classes
attributes
DOM elements
```

not individual CSS properties.

---

# `tst/`

Contains tests.

Tests are intentionally outside `src/`.

Testing philosophy:

> Test useful application behavior, not implementation details.

The goal is functional coverage.

Not line coverage.

Not function coverage.

Not maximizing the number of tests.

Focus primarily on:

```text
happy path
common mistakes
important failures
permissions
validation
state changes
real user flows
```

---

## API tests

API tests should exercise complete HTTP behavior.

They should go through:

```text
request
→ router
→ authentication
→ validation
→ endpoint
→ database
→ response
```

Use real HTTP handlers through `httptest`.

Use a temporary real SQLite database when practical.

Examples:

```text
login succeeds
login fails with wrong password

create user succeeds
create user rejects missing email
create user rejects duplicate email

update user succeeds
update missing user fails

authenticated endpoint works
unauthenticated request is rejected

delete succeeds
forbidden delete is rejected
```

Do not mock internal functions merely to make tests easier.

Tests should protect behavior visible outside the implementation.

---

## UI tests

UI tests should exercise complete browser flows.

Use `chromedp`.

Think like a user.

Examples:

```text
open login
→ enter correct credentials
→ login
→ arrive at application
```

```text
open users
→ create user
→ user appears
```

```text
create user
→ leave email empty
→ submit
→ validation error appears
→ correct email
→ submit
→ succeeds
```

```text
open protected page while logged out
→ redirected to login
→ authenticate
→ continue
```

Prioritize:

### Happy path

Can the user actually complete the important operation?

### Common mistakes

What is the user realistically going to do wrong?

Examples:

```text
empty field
bad value
duplicate value
wrong password
double click
wrong navigation
stale record
missing record
unauthorized action
```

Do not try to enumerate every theoretically possible input.

Do not create a test for every DOM element.

Do not test CSS internals.

Do not test internal JavaScript functions individually by default.

Test the application.

---

## Test helpers

Helpers should remain simple functions.

Example:

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

API helpers:

```go
get()
postJSON()
putJSON()
delete()
login()
expectStatus()
expectJSON()
```

Tests should read almost like scripts:

```go
func TestCreateUserFlow(t *testing.T) {
    loginUI(ctx)

    open(ctx, "/users")
    click(ctx, "#new-user")

    typeText(ctx, "#name", "Daniel")
    typeText(ctx, "#email", "daniel@test.com")

    click(ctx, "#save")

    expectText(ctx, "#users", "Daniel")
}
```

Do not build a test framework.

Do not introduce:

```text
BDD
Gherkin
Page Object architecture
test DSLs
test dependency injection
large fixture systems
```

unless a real problem later requires one.

---

# Unit tests

Do not create unit tests by default.

Internal functions should remain simple enough that testing them through actual application behavior is normally sufficient.

If isolated logic becomes genuinely complicated, a unit test may be added.

That should be exceptional.

Do not create unit tests merely because a function exists.

---

# Runtime assertions

Internal programming invariants may use runtime assertions or explicit failure checks.

Use them when failure indicates a programming bug.

Example concepts:

```text
database must already be initialized
internal ID must already exist here
this function must never receive nil
this state should be impossible
```

Do not use assertions for expected user errors.

These must be handled normally:

```text
invalid request
wrong password
missing resource
expired token
bad input
permission denied
```

The distinction is:

```text
user mistake     → application error

programmer mistake
or impossible state
                 → assertion / loud failure
```

---

# `run/`

Contains executable development and operational scripts.

Scripts should be boring and obvious.

Prefer shell scripts.

Example:

```text
run/dev.sh
run/build.sh
run/test.sh
run/docker.sh
```

Do not create elaborate task runners or build systems unless required.

---

## `run/dev.sh`

Starts the application for local development.

Example behavior:

```bash
#!/usr/bin/env bash
set -e

go run ./src
```

If setup is needed, keep it visible in this script.

Do not hide important development behavior behind tooling.

---

## `run/build.sh`

Builds the application.

Example:

```bash
#!/usr/bin/env bash
set -e

mkdir -p bin

CGO_ENABLED=0 go build \
    -o bin/app \
    ./src
```

Build output belongs in:

```text
bin/
```

---

## `run/test.sh`

Runs all functional tests.

Example:

```bash
#!/usr/bin/env bash
set -e

go test ./tst -v
```

If UI tests require Chrome, the script may verify that Chrome/Chromium is available before running.

Keep this check simple and explicit.

---

## `run/docker.sh`

Builds or runs the Docker image.

Keep Docker operations behind simple commands.

Example:

```bash
#!/usr/bin/env bash
set -e

docker build -t app .
docker run --rm \
    -p 8080:8080 \
    -v "$(pwd)/data:/data" \
    app
```

Do not create deployment abstractions before deployment actually needs them.

---

# `data/`

Contains runtime data and database-related assets.

Example:

```text
data/
    app.db

    sql/
        schema.sql
        seed.sql

        scripts/
            reset.sql
            sample-data.sql
```

Keep SQL visible and editable.

Do not hide schema behavior behind an ORM.

---

# `data/sql/schema.sql`

Contains the current database schema.

Prefer plain SQL.

Example:

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);
```

---

# `data/sql/seed.sql`

Optional development/test seed data.

Keep it deterministic and understandable.

Do not create large fictional datasets unless needed.

---

# `data/sql/scripts/`

Contains direct-purpose SQL scripts.

Examples:

```text
reset.sql
sample-data.sql
cleanup-old-data.sql
repair-user-emails.sql
```

Each script should do one understandable thing.

A SQL script is preferred over building an internal migration or administration framework for a one-off database operation.

---

# Database changes

Keep database evolution simple.

Do not introduce a migration framework automatically.

For a small application it is acceptable to maintain explicit SQL scripts such as:

```text
001_initial.sql
002_add_user_status.sql
003_add_sessions.sql
```

if migrations become necessary.

The goal is not to avoid migrations.

The goal is to avoid adopting migration infrastructure before it is useful.

---

# `doc/`

Contains specifications, decisions and plans.

Documentation must be useful.

No ceremonies.

Do not create documentation merely because a process says a document should exist.

Prefer a few living documents over many templates.

Typical files:

```text
doc/spec.md
doc/plan.md
```

Additional documents may be created when a real feature or decision requires them.

---

## `doc/spec.md`

Defines what the application currently needs to do.

Focus on observable behavior.

Example:

```markdown
# Users

An authenticated administrator can:

- list users,
- create users,
- edit users,
- disable users.

Creating a user requires:

- name,
- unique email.

Common failures:

- missing email,
- invalid email,
- duplicate email,
- unauthorized request.
```

Specifications should be concrete.

Prefer:

```text
user does X
system does Y
```

over abstract architectural language.

---

## `doc/plan.md`

Contains the current implementation plan.

Keep it short and alive.

Example:

```markdown
# Current plan

1. Login endpoint.
2. JWT authentication.
3. Users CRUD.
4. Users UI.
5. API functional tests.
6. UI happy-path test.
7. UI common-error tests.
```

Plans are working notes.

They are not contracts.

They may be changed whenever reality changes.

No project-management ceremony is implied.

---

# Documentation rule

Do not automatically create:

```text
ADR
RFC
architecture document
design review
technical proposal
decision record
meeting notes
status document
```

for every change.

Create documentation only when it helps the user build, understand, operate or modify the application.

A one-line decision may remain one line.

---

# `ai/`

Contains instructions and reusable capabilities for AI-assisted development.

The purpose is to make it possible for one person to build and maintain the application with AI while preserving the project's simplicity.

The AI should behave as an implementation partner, not as an autonomous product manager.

The user is the final decision maker.

---

# `ai/skills/`

Contains reusable task instructions.

Skills should teach the AI how to perform recurring operations in this project.

Examples:

```text
ai/skills/add-api-endpoint.md
ai/skills/add-ui-flow.md
ai/skills/add-functional-test.md
ai/skills/change-database.md
ai/skills/debug.md
```

Skills should be short, procedural, and direct.

Example:

```markdown
# Add API endpoint

1. Read the relevant spec.
2. Inspect existing nearby endpoint code.
3. Add the route.
4. Implement a plain handler function.
5. Call database functions directly.
6. Handle expected validation errors.
7. Add functional API tests:
   - happy path
   - realistic common mistakes
8. Run tests.
9. Update spec only if observable behavior changed.
```

Skills must not introduce new architecture.

---

# `ai/hooks/`

Contains automatic checks or simple AI workflow hooks.

Hooks should protect important project rules.

Examples:

```text
ai/hooks/check-tests.sh
ai/hooks/check-build.sh
ai/hooks/check-simple-architecture.sh
```

Keep hooks small.

Examples of useful checks:

```text
does the project build?
do functional tests pass?
does frontend still run without npm?
were unexpected dependencies added?
were new architectural folders introduced?
```

Do not create complicated AI orchestration.

No autonomous chains of agents unless the user explicitly asks for them.

---

# AI operating rules

The AI must assume that the user is the final owner and final decision maker.

Do not silently make product decisions.

Do not invent requirements.

Do not infer features because they are "industry standard".

Do not add features merely because similar applications normally contain them.

Examples:

Do not automatically add:

```text
password reset
email verification
OAuth
roles
permissions
pagination
audit log
refresh tokens
dark mode
analytics
background jobs
Redis
queues
caching
WebSockets
rate limiting
microservices
```

unless the current requirement needs them.

If something is not specified and a decision is required:

1. Prefer the smallest behavior compatible with the existing code and specification.
2. Make the assumption explicit.
3. Do not silently expand scope.
4. Keep the implementation easy to change.

The AI should not redesign working code merely because another pattern is more fashionable.

---

# AI implementation style

Prefer code that is easy for another AI session to understand immediately.

This means:

```text
plain data
plain structs
plain functions
direct calls
explicit control flow
obvious names
few dependencies
few hidden conventions
```

Avoid clever code.

Avoid excessive generic programming.

Avoid reflection unless genuinely necessary.

Avoid meta-programming.

Avoid implicit dependency graphs.

Avoid abstractions whose purpose cannot be understood from their name and implementation.

The code should be readable by opening the files and following the calls.

---

# No ceremonies

This project explicitly avoids development ceremony that does not produce useful software.

Do not require:

```text
story points
sprints
epics
formal design phases
formal approval phases
RFC processes
mandatory ADRs
code-generation scaffolds
boilerplate architecture
template documents
coverage targets
test-count targets
```

The normal workflow is:

```text
understand requirement
→ inspect current code
→ implement directly
→ test useful behavior
→ run application
→ verify result
```

Documentation and planning exist to support this flow, not control it.

---

# Feature workflow

For a normal feature:

```text
1. Read the relevant spec.

2. Inspect the current implementation.

3. Implement the simplest complete behavior.

4. Use:
   structs
   functions
   direct calls
   plain SQL
   plain JS
   plain CSS

5. Add functional API tests when API behavior changed.

6. Add UI flow tests when user-visible behavior changed.

7. Test:
   happy path
   realistic common mistakes

8. Run:
   run/test.sh

9. Run:
   run/build.sh

10. Update documentation only where reality changed.
```

No additional ceremony is required.

---

# Core project rule

Every new piece of code should make the application do something useful.

Every new abstraction must have an immediate concrete purpose.

Every new dependency must solve an actual problem.

Every new test must protect meaningful behavior.

Every new document must help someone understand or build the application.

Every new folder must have a clear practical reason to exist.

The target is:

**one person, assisted by AI, building a complete application through plain code, simple composition, useful functional tests, and almost no ceremony.**
