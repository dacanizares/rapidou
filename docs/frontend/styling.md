## Frontend philosophy

The frontend must follow the same simplicity rules as the backend.

Use:

* plain HTML,
* plain CSS,
* plain JavaScript,
* native browser APIs.

Do not use a frontend framework.

Do not use a CSS framework by default.

Do not use Bootstrap, Tailwind, component libraries, CSS-in-JS, Sass, Less, PostCSS, or any CSS build process unless a concrete requirement later proves that plain CSS is insufficient.

The frontend must work directly in the browser from the source files embedded in the Go binary.

There should be no frontend compilation step.

## CSS

Use one plain CSS system built with native CSS.

Prefer:

* CSS variables,
* reusable CSS classes,
* Flexbox,
* CSS Grid,
* responsive media queries,
* modern native CSS.

Keep styling out of JavaScript whenever practical.

Do not generate inline styles from JavaScript for normal UI styling.

JavaScript should primarily change:

* content,
* classes,
* attributes,
* visibility,
* state.

Example:

```css
:root {
    --color-bg: #f5f6f8;
    --color-surface: #ffffff;
    --color-text: #202124;
    --color-muted: #687078;

    --color-primary: #2563eb;
    --color-danger: #dc2626;

    --border-color: #d8dde3;

    --radius-sm: 4px;
    --radius-md: 8px;

    --space-1: 4px;
    --space-2: 8px;
    --space-3: 12px;
    --space-4: 16px;
    --space-6: 24px;
    --space-8: 32px;
}
```

Use reusable component classes:

```css
.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;

    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-md);
    border: 1px solid transparent;

    cursor: pointer;
}

.btn-primary {
    background: var(--color-primary);
    color: white;
}

.btn-danger {
    background: var(--color-danger);
    color: white;
}

.card {
    padding: var(--space-4);
    background: var(--color-surface);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
}
```

Use simple layout classes:

```css
.row {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

.stack {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
}

.grid {
    display: grid;
    gap: var(--space-4);
}

.grid-2 {
    grid-template-columns: repeat(2, minmax(0, 1fr));
}

.grid-3 {
    grid-template-columns: repeat(3, minmax(0, 1fr));
}
```

And application-specific classes when they make the code clearer:

```css
.section-users {
    display: grid;
    gap: var(--space-4);
}

.users-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.user-row {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: center;
    gap: var(--space-3);
}
```

This is preferred over large collections of microscopic utility classes.

HTML should remain readable:

```html
<section class="section-users">
    <div class="users-toolbar">
        <h1>Users</h1>

        <button class="btn btn-primary" id="new-user">
            New user
        </button>
    </div>

    <div class="card" id="users"></div>
</section>
```

The meaning of the UI should be visible directly from the HTML.

## CSS organization

Start with a single file:

```text
web/
    index.html
    app.js
    app.css
```

Organize `app.css` in simple sections:

```css
/* variables */
/* reset */
/* base */
/* layout */
/* components */
/* forms */
/* application sections */
/* state */
/* responsive */
```

Do not split CSS into many files merely for organization.

Split only when the file becomes genuinely difficult to navigate.

## State classes

Use simple classes for visual state:

```css
.hidden {
    display: none;
}

.is-loading {
    opacity: 0.6;
    pointer-events: none;
}

.is-error {
    border-color: var(--color-danger);
}
```

JavaScript may simply add or remove them:

```js
element.classList.add("is-loading");

element.classList.remove("is-loading");
```

Prefer this over JavaScript manually manipulating individual CSS properties.

## Rendering

Rendering functions should update only the relevant part of the interface when practical.

Example:

```js
async function refreshUsers() {
    state.users = await api("/api/users");
    renderUsers("#users");
}
```

And:

```js
function renderUsers(selector) {
    const root = document.querySelector(selector);

    root.replaceChildren(
        ...state.users.map(UserRow)
    );
}
```

There is no requirement for a global application render cycle.

If only users changed, render users.

If only navigation changed, render navigation.

If only a dialog changed, update the dialog.

Keep updates local and explicit.

## UI components

Components are simply:

* HTML structures,
* CSS classes,
* JavaScript functions.

Example:

```js
function UserRow(user) {
    const row = document.createElement("div");
    row.className = "user-row";

    const name = document.createElement("span");
    name.textContent = user.name;

    const remove = document.createElement("button");
    remove.className = "btn btn-danger";
    remove.textContent = "Delete";

    remove.onclick = () => deleteUser(user.id);

    row.append(name, remove);

    return row;
}
```

Do not introduce a component framework to represent this concept.

## Design rule

Create a small visual language for the application using CSS classes.

For example:

```text
.btn
.btn-primary
.btn-danger

.card

.input
.input-error

.row
.stack
.grid
.grid-2

.dialog

.table

.badge
.badge-success
.badge-error

.section-users
.section-settings
```

Keep this vocabulary small.

Add a CSS class when a real visual concept appears.

Do not attempt to design a complete design system in advance.

The target is:

**modern UI built from HTML, CSS classes, browser APIs, and simple JavaScript functions — with no frontend machinery between the source code and the browser.**
