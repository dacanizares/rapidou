package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const sessionCookie = "rapidou_session"

type claims struct {
	Subject int64 `json:"sub"`
	Expires int64 `json:"exp"`
}

type userContextKey struct{}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input loginInput
	if !readJSON(w, r, &input) {
		return
	}
	user, hash, err := userCredentialsByEmail(app.DB, strings.TrimSpace(input.Email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, err := signToken(user.ID, time.Now().Add(12*time.Hour), app.JWTSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: token, Path: "/", HttpOnly: true,
		Secure: app.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: 12 * 60 * 60,
	})
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (app *App) logoutHandler(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: "", Path: "/", HttpOnly: true,
		Secure: app.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: -1,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) meHandler(w http.ResponseWriter, r *http.Request) {
	user, err := userByID(app.DB, currentUserID(r))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""
		if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		} else if cookie, err := r.Cookie(sessionCookie); err == nil {
			token = cookie.Value
		}
		userID, err := verifyToken(token, app.JWTSecret, time.Now())
		if err != nil {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey{}, userID)))
	})
}

func currentUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value(userContextKey{}).(int64)
	if !ok || userID == 0 {
		panic("authenticated request has no user ID")
	}
	return userID
}

func signToken(userID int64, expires time.Time, secret []byte) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadJSON, err := json.Marshal(claims{Subject: userID, Expires: expires.Unix()})
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsigned := header + "." + payload
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(unsigned))
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func verifyToken(token string, secret []byte, now time.Time) (int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, errors.New("invalid token")
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(signature, mac.Sum(nil)) {
		return 0, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, errors.New("invalid payload")
	}
	var value claims
	if json.Unmarshal(payload, &value) != nil || value.Subject <= 0 || value.Expires <= now.Unix() {
		return 0, errors.New("invalid claims")
	}
	return value.Subject, nil
}

func parseID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}
