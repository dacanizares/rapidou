package tst

import (
	"net/http"
	"testing"
)

type user struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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
}
