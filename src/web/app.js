const state = { currentUser: null, games: [] };

async function api(path, options = {}) {
    const response = await fetch(path, {
        ...options,
        headers: { "Content-Type": "application/json", ...(options.headers || {}) }
    });
    const body = response.status === 204 ? null : await response.json().catch(() => null);
    if (!response.ok) throw new Error(body?.error || `La solicitud falló (${response.status})`);
    return body;
}

function showAuthenticated(authenticated) {
    document.querySelector("#login-view").classList.toggle("hidden", authenticated);
    document.querySelector("#curator-view").classList.toggle("hidden", !authenticated);
    document.querySelector("#logout").classList.toggle("hidden", !authenticated);
}

function GameCard(game) {
    const article = document.createElement("article");
    article.className = "artifact card";

    const gallery = document.createElement("div");
    gallery.className = "artifact-gallery";
    if (game.images.length === 0) {
        const placeholder = document.createElement("div");
        placeholder.className = "artifact-placeholder";
        placeholder.textContent = "Sin imagen";
        gallery.append(placeholder);
    } else {
        gallery.append(...game.images.map((image) => {
            const element = document.createElement("img");
            element.src = image.src;
            element.alt = image.alt;
            element.loading = "lazy";
            return element;
        }));
    }

    const year = document.createElement("span");
    year.className = "artifact-year";
    year.textContent = game.release_year;

    const title = document.createElement("h3");
    title.textContent = game.title;

    const platform = document.createElement("p");
    platform.className = "artifact-platform";
    platform.textContent = game.platform;

    const description = document.createElement("p");
    description.className = "artifact-description";
    description.textContent = game.description || "Pieza sin descripción catalogada.";

    article.append(gallery, year, title, platform, description);
    if (state.currentUser) {
        const upload = document.createElement("form");
        upload.className = "image-upload";
        const input = document.createElement("input");
        input.type = "file";
        input.name = "image";
        input.accept = "image/jpeg,image/png,image/gif,image/webp";
        input.required = true;
        input.setAttribute("aria-label", `Imagen para ${game.title}`);
        const uploadButton = document.createElement("button");
        uploadButton.className = "btn";
        uploadButton.type = "submit";
        uploadButton.textContent = "Subir imagen";
        const uploadMessage = document.createElement("span");
        uploadMessage.className = "upload-message";
        upload.addEventListener("submit", (event) => uploadGameImage(game, event, uploadMessage));
        upload.append(input, uploadButton, uploadMessage);

        const remove = document.createElement("button");
        remove.className = "btn btn-danger";
        remove.type = "button";
        remove.textContent = "Retirar pieza";
        remove.addEventListener("click", () => deleteGame(game));
        article.append(upload, remove);
    }
    return article;
}

function renderGames() {
    document.querySelector("#games").replaceChildren(...state.games.map(GameCard));
    document.querySelector("#empty-collection").classList.toggle("hidden", state.games.length > 0);
    document.querySelector("#game-count").textContent = `${state.games.length} ${state.games.length === 1 ? "pieza" : "piezas"}`;
}

async function refreshGames() {
    state.games = await api("/api/games");
    renderGames();
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
        document.querySelector("#current-user").textContent = `Curador: ${result.user.name}`;
        showAuthenticated(true);
        renderGames();
    } catch (cause) {
        error.textContent = cause.message;
    }
}

async function createGame(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const error = document.querySelector("#game-error");
    error.textContent = "";
    try {
        await api("/api/games", {
            method: "POST",
            body: JSON.stringify({
                title: form.elements.title.value,
                platform: form.elements.platform.value,
                release_year: Number(form.elements.release_year.value),
                description: form.elements.description.value
            })
        });
        form.reset();
        await refreshGames();
    } catch (cause) {
        error.textContent = cause.message;
    }
}

async function deleteGame(game) {
    if (!window.confirm(`¿Retirar ${game.title} de la colección?`)) return;
    await api(`/api/games/${game.id}`, { method: "DELETE" });
    await refreshGames();
}

async function uploadGameImage(game, event, message) {
    event.preventDefault();
    message.textContent = "";
    const form = event.currentTarget;
    const data = new FormData(form);
    data.set("alt", game.title);
    try {
        const response = await fetch(`/api/games/${game.id}/images`, { method: "POST", body: data });
        const body = await response.json().catch(() => null);
        if (!response.ok) throw new Error(body?.error || "No se pudo subir la imagen");
        form.reset();
        await refreshGames();
    } catch (cause) {
        message.textContent = cause.message;
    }
}

async function logout() {
    await api("/api/logout", { method: "POST" });
    state.currentUser = null;
    document.querySelector("#current-user").textContent = "";
    showAuthenticated(false);
    renderGames();
}

async function start() {
    document.querySelector("#login-form").addEventListener("submit", login);
    document.querySelector("#game-form").addEventListener("submit", createGame);
    document.querySelector("#logout").addEventListener("click", logout);
    await refreshGames();
    try {
        state.currentUser = await api("/api/me");
        document.querySelector("#current-user").textContent = `Curador: ${state.currentUser.name}`;
        showAuthenticated(true);
        renderGames();
    } catch {
        showAuthenticated(false);
    }
}

start();
