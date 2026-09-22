// Package tst verifies an application exclusively through its HTTP boundary.
package tst

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Application = http.Handler
type Factory func(t *testing.T) (Application, func())

type Credentials struct {
	Email    string
	Password string
}

type response struct {
	Status int
	Header http.Header
	Body   []byte
}

func request(t *testing.T, app http.Handler, method, path, token string, body any) response {
	t.Helper()
	var content io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		content = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, path, content)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, req)
	result := recorder.Result()
	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	return response{Status: result.StatusCode, Header: result.Header, Body: data}
}

func expectStatus(t *testing.T, result response, expected int) {
	t.Helper()
	if result.Status != expected {
		t.Fatalf("expected HTTP %d, got %d: %s", expected, result.Status, result.Body)
	}
}

func decode[T any](t *testing.T, result response) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(result.Body, &value); err != nil {
		t.Fatalf("decode response %q: %v", result.Body, err)
	}
	return value
}

func login(t *testing.T, app http.Handler, credentials Credentials) string {
	t.Helper()
	result := request(t, app, http.MethodPost, "/api/login", "", map[string]string{
		"email": credentials.Email, "password": credentials.Password,
	})
	expectStatus(t, result, http.StatusOK)
	value := decode[struct {
		Token string `json:"token"`
	}](t, result)
	if value.Token == "" {
		t.Fatal("login response contains no token")
	}
	return value.Token
}
