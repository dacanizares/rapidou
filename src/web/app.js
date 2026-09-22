const state = { currentUser: null, users: [] };

async function api(path, options = {}) {
    const response = await fetch(path, {
        ...options,
        headers: { "Content-Type": "application/json", ...(options.headers || {}) }
    });
    const body = response.status === 204 ? null : await response.json().catch(() => null);
    if (!response.ok) throw new Error(body?.error || `Request failed (${response.status})`);
    return body;
}

function showAuthenticated(authenticated) {
    document.querySelector("#login-view").classList.toggle("hidden", authenticated);
    document.querySelector("#app-view").classList.toggle("hidden", !authenticated);
}

function renderUsers() {
    const root = document.querySelector("#users");
    root.replaceChildren(...state.users.map((user) => {
        const row = document.createElement("div");
        row.className = "user-row";
        const identity = document.createElement("div");
        const name = document.createElement("strong");
        name.textContent = user.name;
        const email = document.createElement("span");
        email.className = "muted";
        email.textContent = user.email;
        identity.append(name, email);
        const remove = document.createElement("button");
        remove.className = "btn btn-danger";
        remove.type = "button";
        remove.textContent = "Delete";
        remove.disabled = user.id === state.currentUser.id;
        remove.addEventListener("click", () => deleteUser(user));
        row.append(identity, remove);
        return row;
    }));
}

async function refreshUsers() {
    state.users = await api("/api/users");
    renderUsers();
}

async function login(event) {
    event.preventDefault();
    const error = document.querySelector("#login-error");
    error.textContent = "";
    try {
        const result = await api("/api/login", {
            method: "POST",
            body: JSON.stringify({
                email: document.querySelector("#login-email").value,
                password: document.querySelector("#login-password").value
            })
        });
        state.currentUser = result.user;
        document.querySelector("#current-user").textContent = result.user.name;
        showAuthenticated(true);
        await refreshUsers();
    } catch (cause) {
        error.textContent = cause.message;
    }
}

async function createUser(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const error = document.querySelector("#user-error");
    error.textContent = "";
    try {
        await api("/api/users", {
            method: "POST",
            body: JSON.stringify({
                name: form.elements.name.value,
                email: form.elements.email.value,
                password: form.elements.password.value
            })
        });
        form.reset();
        await refreshUsers();
    } catch (cause) {
        error.textContent = cause.message;
    }
}

async function deleteUser(user) {
    if (!window.confirm(`Delete ${user.name}?`)) return;
    await api(`/api/users/${user.id}`, { method: "DELETE" });
    await refreshUsers();
}

async function logout() {
    await api("/api/logout", { method: "POST" });
    state.currentUser = null;
    state.users = [];
    showAuthenticated(false);
}

async function start() {
    document.querySelector("#login-form").addEventListener("submit", login);
    document.querySelector("#user-form").addEventListener("submit", createUser);
    document.querySelector("#logout").addEventListener("click", logout);
    try {
        state.currentUser = await api("/api/me");
        document.querySelector("#current-user").textContent = state.currentUser.name;
        showAuthenticated(true);
        await refreshUsers();
    } catch {
        showAuthenticated(false);
    }
}

start();

