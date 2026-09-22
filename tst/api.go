package tst

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type user struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type game struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Platform    string `json:"platform"`
	ReleaseYear int    `json:"release_year"`
	SteamURL    string `json:"steam_url"`
	GOGURL      string `json:"gog_url"`
	Images      []struct {
		ID   int64  `json:"id"`
		Src  string `json:"src"`
		Kind string `json:"kind"`
	} `json:"images"`
}

// RunAPI executes the portable HTTP contract against a fresh app per flow.
func RunAPI(t *testing.T, factory Factory, credentials Credentials) {
	t.Helper()

	t.Run("health and authentication", func(t *testing.T) {
		app, close := factory(t)
		defer close()

		expectStatus(t, request(t, app, http.MethodGet, "/api/health", "", nil), http.StatusOK)
		expectStatus(t, request(t, app, http.MethodGet, "/api/users", "", nil), http.StatusUnauthorized)
		expectStatus(t, request(t, app, http.MethodPost, "/api/login", "", map[string]string{
			"email": credentials.Email, "password": "wrong-password",
		}), http.StatusUnauthorized)

		token := login(t, app, credentials)
		expectStatus(t, request(t, app, http.MethodGet, "/api/me", token, nil), http.StatusOK)
	})

	t.Run("user lifecycle", func(t *testing.T) {
		app, close := factory(t)
		defer close()
		token := login(t, app, credentials)

		created := request(t, app, http.MethodPost, "/api/users", token, map[string]string{
			"name": "Daniel", "email": "daniel@rapidou.test", "password": "long-enough",
		})
		expectStatus(t, created, http.StatusCreated)
		newUser := decode[user](t, created)
		if newUser.ID == 0 || newUser.Email != "daniel@rapidou.test" {
			t.Fatalf("unexpected created user: %+v", newUser)
		}

		listed := request(t, app, http.MethodGet, "/api/users", token, nil)
		expectStatus(t, listed, http.StatusOK)
		users := decode[[]user](t, listed)
		if len(users) != 2 {
			t.Fatalf("expected initial and created users, got %+v", users)
		}

		updated := request(t, app, http.MethodPut, "/api/users/2", token, map[string]string{
			"name": "Daniel C", "email": "daniel@rapidou.test", "password": "",
		})
		expectStatus(t, updated, http.StatusOK)
		if got := decode[user](t, updated).Name; got != "Daniel C" {
			t.Fatalf("expected updated name, got %q", got)
		}

		expectStatus(t, request(t, app, http.MethodDelete, "/api/users/2", token, nil), http.StatusNoContent)
		expectStatus(t, request(t, app, http.MethodDelete, "/api/users/2", token, nil), http.StatusNotFound)
	})

	t.Run("important user errors", func(t *testing.T) {
		app, close := factory(t)
		defer close()
		token := login(t, app, credentials)

		expectStatus(t, request(t, app, http.MethodPost, "/api/users", token, map[string]string{
			"name": "Missing", "email": "", "password": "long-enough",
		}), http.StatusBadRequest)
		expectStatus(t, request(t, app, http.MethodPost, "/api/users", token, map[string]string{
			"name": "Duplicate", "email": credentials.Email, "password": "long-enough",
		}), http.StatusConflict)
		expectStatus(t, request(t, app, http.MethodPut, "/api/users/9999", token, map[string]string{
			"name": "Missing", "email": "missing@rapidou.test", "password": "",
		}), http.StatusNotFound)
		expectStatus(t, request(t, app, http.MethodDelete, "/api/users/1", token, nil), http.StatusConflict)
	})

	t.Run("museum piece lifecycle", func(t *testing.T) {
		app, close := factory(t)
		defer close()

		listed := request(t, app, http.MethodGet, "/api/games", "", nil)
		expectStatus(t, listed, http.StatusOK)
		if games := decode[[]game](t, listed); len(games) != 0 {
			t.Fatalf("expected an empty public museum, got %+v", games)
		}

		piece := map[string]any{
			"title": "Super Mario Bros.", "platform": "NES", "release_year": 1985,
			"description": "A defining side-scrolling platform game.",
			"steam_url":   "https://store.steampowered.com/app/123", "gog_url": "https://www.gog.com/en/game/example",
		}
		expectStatus(t, request(t, app, http.MethodPost, "/api/games", "", piece), http.StatusUnauthorized)
		token := login(t, app, credentials)

		created := request(t, app, http.MethodPost, "/api/games", token, piece)
		expectStatus(t, created, http.StatusCreated)
		newGame := decode[game](t, created)
		if newGame.ID == 0 || newGame.Title != "Super Mario Bros." || newGame.SteamURL == "" || newGame.GOGURL == "" {
			t.Fatalf("unexpected museum piece: %+v", newGame)
		}

		imageData := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}
		uploaded := uploadImage(t, app, fmt.Sprintf("/api/games/%d/images", newGame.ID), token, "museum.png", imageData)
		expectStatus(t, uploaded, http.StatusCreated)
		image := decode[struct {
			ID  int64  `json:"id"`
			Src string `json:"src"`
		}](t, uploaded)
		if image.ID == 0 || image.Src == "" {
			t.Fatalf("unexpected uploaded image: %+v", image)
		}
		served := request(t, app, http.MethodGet, image.Src, "", nil)
		expectStatus(t, served, http.StatusOK)
		if served.Header.Get("Content-Type") != "image/png" {
			t.Fatalf("expected uploaded PNG, got %q", served.Header.Get("Content-Type"))
		}
		linked := request(t, app, http.MethodPost, fmt.Sprintf("/api/games/%d/image-links", newGame.ID), token, map[string]string{
			"url": "https://example.com/museum.jpg", "alt": "Remote museum image",
		})
		expectStatus(t, linked, http.StatusCreated)
		linkedImage := decode[struct {
			ID   int64  `json:"id"`
			Kind string `json:"kind"`
		}](t, linked)
		if linkedImage.Kind != "url" {
			t.Fatalf("expected URL image kind, got %+v", linkedImage)
		}

		updated := request(t, app, http.MethodPut, "/api/games/1", token, map[string]any{
			"title": "Super Mario Bros.", "platform": "Nintendo Entertainment System", "release_year": 1985,
			"description": "A defining side-scrolling platform game.",
			"steam_url":   "https://store.steampowered.com/app/456", "gog_url": "",
		})
		expectStatus(t, updated, http.StatusOK)
		if got := decode[game](t, updated); got.Platform != "Nintendo Entertainment System" || got.SteamURL != "https://store.steampowered.com/app/456" || got.GOGURL != "" {
			t.Fatalf("expected updated platform and store links, got %+v", got)
		}

		listed = request(t, app, http.MethodGet, "/api/games", "", nil)
		expectStatus(t, listed, http.StatusOK)
		if games := decode[[]game](t, listed); len(games) != 1 {
			t.Fatalf("expected one public museum piece, got %+v", games)
		} else if len(games[0].Images) != 2 {
			t.Fatalf("expected uploaded and URL images in public museum, got %+v", games[0].Images)
		}
		expectStatus(t, request(t, app, http.MethodDelete, fmt.Sprintf("/api/game-images/%d", image.ID), token, nil), http.StatusNoContent)
		expectStatus(t, request(t, app, http.MethodDelete, fmt.Sprintf("/api/game-images/%d", linkedImage.ID), token, nil), http.StatusNoContent)

		expectStatus(t, request(t, app, http.MethodDelete, "/api/games/1", token, nil), http.StatusNoContent)
		expectStatus(t, request(t, app, http.MethodDelete, "/api/games/1", token, nil), http.StatusNotFound)
	})

	t.Run("museum piece validation", func(t *testing.T) {
		app, close := factory(t)
		defer close()
		token := login(t, app, credentials)

		expectStatus(t, request(t, app, http.MethodPost, "/api/games", token, map[string]any{
			"title": "", "platform": "Arcade", "release_year": 1980, "description": "",
		}), http.StatusBadRequest)
		expectStatus(t, request(t, app, http.MethodPost, "/api/games", token, map[string]any{
			"title": "Future artifact", "platform": "Unknown", "release_year": 2200, "description": "",
		}), http.StatusBadRequest)
		expectStatus(t, request(t, app, http.MethodPost, "/api/games", token, map[string]any{
			"title": "Bad store", "platform": "PC", "release_year": 2000, "description": "",
			"steam_url": "https://example.com/not-steam", "gog_url": "http://gog.com/not-secure",
		}), http.StatusBadRequest)
		expectStatus(t, request(t, app, http.MethodPut, "/api/games/9999", token, map[string]any{
			"title": "Missing", "platform": "Arcade", "release_year": 1980, "description": "",
		}), http.StatusNotFound)
		expectStatus(t, uploadImage(t, app, "/api/games/9999/images", token, "note.txt", []byte("not an image")), http.StatusNotFound)
		created := request(t, app, http.MethodPost, "/api/games", token, map[string]any{
			"title": "Image test", "platform": "PC", "release_year": 2000, "description": "",
		})
		expectStatus(t, created, http.StatusCreated)
		game := decode[game](t, created)
		expectStatus(t, uploadImage(t, app, fmt.Sprintf("/api/games/%d/images", game.ID), token, "note.txt", []byte("not an image")), http.StatusBadRequest)
		expectStatus(t, request(t, app, http.MethodPost, fmt.Sprintf("/api/games/%d/image-links", game.ID), token, map[string]string{
			"url": "http://example.com/not-secure.jpg", "alt": "Invalid",
		}), http.StatusBadRequest)
	})
}

// RunSampleMuseum verifies the development seed through the public API.
func RunSampleMuseum(t *testing.T, factory Factory) {
	t.Helper()
	app, close := factory(t)
	defer close()
	listed := request(t, app, http.MethodGet, "/api/games", "", nil)
	expectStatus(t, listed, http.StatusOK)
	games := decode[[]game](t, listed)
	if len(games) != 6 {
		t.Fatalf("expected six seeded museum pieces, got %d", len(games))
	}
	for _, game := range games {
		if !strings.HasPrefix(game.SteamURL, "https://store.steampowered.com/app/") {
			t.Fatalf("expected Steam seed URL for %s, got %q", game.Title, game.SteamURL)
		}
		if len(game.Images) != 2 {
			t.Fatalf("expected two seed images for %s, got %+v", game.Title, game.Images)
		}
		for _, image := range game.Images {
			if !strings.HasPrefix(image.Src, "https://") {
				t.Fatalf("expected remote HTTPS seed image for %s, got %q", game.Title, image.Src)
			}
			if image.Kind != "url" {
				t.Fatalf("expected URL seed image for %s, got %q", game.Title, image.Kind)
			}
		}
	}
}
