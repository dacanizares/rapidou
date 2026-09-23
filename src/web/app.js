const state = { currentUser: null, games: [], editingGameID: null };

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
    document.querySelector("#open-login").classList.toggle("hidden", authenticated);
    document.querySelector("#new-game").classList.toggle("hidden", !authenticated);
    document.querySelector("#logout").classList.toggle("hidden", !authenticated);
}

function showPopup(dialog) {
    dialog.showModal();
    requestAnimationFrame(() => requestAnimationFrame(() => dialog.classList.add("is-visible")));
}

function closePopup(dialog) {
    dialog.classList.remove("is-visible");
    dialog.close();
}

function StoreLinks(game) {
    const stores = document.createElement("div");
    stores.className = "store-links";
    [["Steam", game.steam_url], ["GOG", game.gog_url]].forEach(([label, href]) => {
        if (!href) return;
        const link = document.createElement("a");
        link.className = `store-link store-${label.toLowerCase()}`;
        link.href = href;
        link.target = "_blank";
        link.rel = "noreferrer";
        link.textContent = label;
        stores.append(link);
    });
    return stores;
}

function openGameDetail(game) {
    document.querySelector("#detail-title").textContent = game.title;
    document.querySelector("#detail-metadata").textContent = `${game.platform} · ${game.release_year}`;
    document.querySelector("#detail-description").textContent = game.description || "Pieza sin descripción catalogada.";

    const gallery = document.querySelector("#detail-gallery");
    if (game.images.length === 0) {
        const placeholder = document.createElement("div");
        placeholder.className = "artifact-placeholder detail-placeholder";
        placeholder.textContent = "Sin imagen";
        gallery.replaceChildren(placeholder);
    } else {
        gallery.replaceChildren(...game.images.map((image) => {
            const element = document.createElement("img");
            element.src = image.src;
            element.alt = image.alt;
            return element;
        }));
    }

    const stores = StoreLinks(game);
    stores.id = "detail-stores";
    stores.classList.add("detail-store-links");
    document.querySelector("#detail-stores").replaceWith(stores);
    showPopup(document.querySelector("#detail-dialog"));
}

function GameCard(game) {
    const article = document.createElement("article");
    article.className = "artifact card";
    article.tabIndex = 0;
    article.setAttribute("aria-label", `Ver ficha de ${game.title}`);
    article.addEventListener("click", (event) => {
        if (!event.target.closest("a, button")) openGameDetail(game);
    });
    article.addEventListener("keydown", (event) => {
        if (event.target !== article || (event.key !== "Enter" && event.key !== " ")) return;
        event.preventDefault();
        openGameDetail(game);
    });

    const gallery = document.createElement("div");
    gallery.className = "artifact-gallery";
    if (game.images.length === 0) {
        const placeholder = document.createElement("div");
        placeholder.className = "artifact-placeholder";
        placeholder.textContent = "Sin imagen";
        gallery.append(placeholder);
    } else {
        gallery.append(...game.images.slice(0, 2).map((image) => {
            const element = document.createElement("img");
            element.src = image.src;
            element.alt = image.alt;
            element.loading = "lazy";
            return element;
        }));
        if (game.images.length > 2) {
            const count = document.createElement("span");
            count.className = "image-count";
            count.textContent = `+${game.images.length - 2}`;
            gallery.append(count);
        }
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
    const stores = StoreLinks(game);
    if (stores.childElementCount > 0) article.append(stores);
    if (state.currentUser) {
        const actions = document.createElement("div");
        actions.className = "artifact-actions";
        const edit = document.createElement("button");
        edit.className = "btn";
        edit.type = "button";
        edit.textContent = "Editar";
        edit.addEventListener("click", () => openEditGame(game));
        const remove = document.createElement("button");
        remove.className = "btn btn-danger";
        remove.type = "button";
        remove.textContent = "Retirar";
        remove.addEventListener("click", () => deleteGame(game));
        actions.append(edit, remove);
        article.append(actions);
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
        closePopup(document.querySelector("#login-dialog"));
    } catch (cause) {
        error.textContent = cause.message;
    }
}

function openNewGame() {
    state.editingGameID = null;
    document.querySelector("#game-form").reset();
    document.querySelector("#dialog-title").textContent = "Agregar pieza";
    document.querySelector("#save").textContent = "Crear pieza";
    document.querySelector("#game-error").textContent = "";
    document.querySelector("#image-manager").classList.add("hidden");
    showPopup(document.querySelector("#game-dialog"));
}

function openEditGame(game) {
    state.editingGameID = game.id;
    const form = document.querySelector("#game-form");
    form.elements.title.value = game.title;
    form.elements.platform.value = game.platform;
    form.elements.release_year.value = game.release_year;
    form.elements.description.value = game.description;
    form.elements.steam_url.value = game.steam_url || "";
    form.elements.gog_url.value = game.gog_url || "";
    document.querySelector("#dialog-title").textContent = `Editar ${game.title}`;
    document.querySelector("#save").textContent = "Guardar cambios";
    document.querySelector("#game-error").textContent = "";
    document.querySelector("#image-error").textContent = "";
    document.querySelector("#image-manager").classList.remove("hidden");
    renderManagedImages(game);
    showPopup(document.querySelector("#game-dialog"));
}

function closeGameDialog() {
    closePopup(document.querySelector("#game-dialog"));
    state.editingGameID = null;
}

async function saveGame(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const error = document.querySelector("#game-error");
    error.textContent = "";
    const input = {
        title: form.elements.title.value,
        platform: form.elements.platform.value,
        release_year: Number(form.elements.release_year.value),
        description: form.elements.description.value,
        steam_url: form.elements.steam_url.value,
        gog_url: form.elements.gog_url.value
    };
    try {
        const creating = state.editingGameID === null;
        const saved = await api(creating ? "/api/games" : `/api/games/${state.editingGameID}`, {
            method: creating ? "POST" : "PUT",
            body: JSON.stringify(input)
        });
        state.editingGameID = saved.id;
        await refreshGames();
        if (!creating) {
            closeGameDialog();
            return;
        }
        const game = state.games.find((item) => item.id === saved.id);
        document.querySelector("#dialog-title").textContent = `Editar ${game.title}`;
        document.querySelector("#save").textContent = "Guardar cambios";
        document.querySelector("#image-manager").classList.remove("hidden");
        renderManagedImages(game);
    } catch (cause) {
        error.textContent = cause.message;
    }
}

function renderManagedImages(game) {
    const root = document.querySelector("#managed-images");
    if (game.images.length === 0) {
        const empty = document.createElement("p");
        empty.className = "manager-empty";
        empty.textContent = "Esta pieza todavía no tiene imágenes.";
        root.replaceChildren(empty);
        return;
    }
    root.replaceChildren(...game.images.map((image) => {
        const row = document.createElement("div");
        row.className = "managed-image";
        const preview = document.createElement("img");
        preview.src = image.src;
        preview.alt = image.alt;
        const details = document.createElement("div");
        details.className = "managed-image-details";
        const kind = document.createElement("span");
        kind.className = `source-badge source-${image.kind}`;
        kind.textContent = image.kind === "url" ? "URL" : "Archivo";
        const source = document.createElement(image.kind === "url" ? "a" : "span");
        source.className = "image-source";
        source.textContent = image.kind === "url" ? image.src : "Guardada en SQLite";
        if (image.kind === "url") {
            source.href = image.src;
            source.target = "_blank";
            source.rel = "noreferrer";
        }
        details.append(kind, source);
        const remove = document.createElement("button");
        remove.className = "btn image-remove";
        remove.type = "button";
        remove.textContent = "Eliminar";
        remove.addEventListener("click", () => deleteGameImage(image.id));
        row.append(preview, details, remove);
        return row;
    }));
}

async function refreshEditor() {
    await refreshGames();
    const game = state.games.find((item) => item.id === state.editingGameID);
    if (game) renderManagedImages(game);
}

async function addImageURL(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const error = document.querySelector("#image-error");
    error.textContent = "";
    try {
        await api(`/api/games/${state.editingGameID}/image-links`, {
            method: "POST",
            body: JSON.stringify({ url: form.elements.url.value, alt: form.elements.alt.value })
        });
        form.reset();
        await refreshEditor();
    } catch (cause) {
        error.textContent = cause.message;
    }
}

async function uploadGameImage(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const error = document.querySelector("#image-error");
    error.textContent = "";
    const data = new FormData(form);
    try {
        const response = await fetch(`/api/games/${state.editingGameID}/images`, { method: "POST", body: data });
        const body = await response.json().catch(() => null);
        if (!response.ok) throw new Error(body?.error || "No se pudo subir la imagen");
        form.reset();
        await refreshEditor();
    } catch (cause) {
        error.textContent = cause.message;
    }
}

async function deleteGameImage(imageID) {
    await api(`/api/game-images/${imageID}`, { method: "DELETE" });
    await refreshEditor();
}

async function deleteGame(game) {
    if (!window.confirm(`¿Retirar ${game.title} de la colección?`)) return;
    await api(`/api/games/${game.id}`, { method: "DELETE" });
    await refreshGames();
}

async function logout() {
    await api("/api/logout", { method: "POST" });
    if (document.querySelector("#game-dialog").open) closeGameDialog();
    state.currentUser = null;
    document.querySelector("#current-user").textContent = "";
    showAuthenticated(false);
    renderGames();
}

async function start() {
    document.querySelector("#login-form").addEventListener("submit", login);
    document.querySelector("#new-game").addEventListener("click", openNewGame);
    document.querySelector("#game-form").addEventListener("submit", saveGame);
    document.querySelector("#image-url-form").addEventListener("submit", addImageURL);
    document.querySelector("#image-upload-form").addEventListener("submit", uploadGameImage);
    document.querySelector("#close-dialog").addEventListener("click", closeGameDialog);
    document.querySelector("#cancel-dialog").addEventListener("click", closeGameDialog);
    document.querySelector("#logout").addEventListener("click", logout);
    document.querySelector("#close-detail").addEventListener("click", () => closePopup(document.querySelector("#detail-dialog")));
    document.querySelector("#detail-dialog").addEventListener("click", (event) => {
        if (event.target === event.currentTarget) closePopup(event.currentTarget);
    });
    document.querySelector("#detail-dialog").addEventListener("close", (event) => event.currentTarget.classList.remove("is-visible"));
    document.querySelector("#open-login").addEventListener("click", () => {
        document.querySelector("#login-error").textContent = "";
        showPopup(document.querySelector("#login-dialog"));
        document.querySelector("#login-email").focus();
    });
    document.querySelector("#close-login").addEventListener("click", () => closePopup(document.querySelector("#login-dialog")));
    document.querySelector("#login-dialog").addEventListener("click", (event) => {
        if (event.target === event.currentTarget) closePopup(event.currentTarget);
    });
    document.querySelector("#login-dialog").addEventListener("close", (event) => event.currentTarget.classList.remove("is-visible"));
    document.querySelector("#game-dialog").addEventListener("click", (event) => {
        if (event.target === event.currentTarget) closeGameDialog();
    });
    document.querySelector("#game-dialog").addEventListener("close", (event) => event.currentTarget.classList.remove("is-visible"));
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
