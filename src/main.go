package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed web/*
var webFiles embed.FS

type Config struct {
	DatabasePath  string
	JWTSecret     string
	CookieSecure  bool
	AdminEmail    string
	AdminPassword string
}

type App struct {
	DB           *sql.DB
	JWTSecret    []byte
	CookieSecure bool
}

func main() {
	address := env("APP_ADDRESS", ":8080")
	config := Config{
		DatabasePath:  env("APP_DATABASE", "data/app.db"),
		JWTSecret:     os.Getenv("APP_JWT_SECRET"),
		CookieSecure:  strings.EqualFold(os.Getenv("APP_COOKIE_SECURE"), "true"),
		AdminEmail:    os.Getenv("APP_ADMIN_EMAIL"),
		AdminPassword: os.Getenv("APP_ADMIN_PASSWORD"),
	}

	if config.JWTSecret == "" {
		log.Fatal("APP_JWT_SECRET is required")
	}

	app, err := newApp(config)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	server := &http.Server{
		Addr:              address,
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("rapidou listening on %s", address)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func newApp(config Config) (*App, error) {
	if config.DatabasePath == "" {
		return nil, errors.New("database path is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWT secret is required")
	}
	if (config.AdminEmail == "") != (config.AdminPassword == "") {
		return nil, errors.New("admin email and password must be provided together")
	}

	if config.DatabasePath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(config.DatabasePath), 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	db, err := openDatabase(config.DatabasePath)
	if err != nil {
		return nil, err
	}

	app := &App{DB: db, JWTSecret: []byte(config.JWTSecret), CookieSecure: config.CookieSecure}
	if config.AdminEmail != "" {
		if err := ensureInitialUser(db, config.AdminEmail, config.AdminPassword); err != nil {
			db.Close()
			return nil, fmt.Errorf("create initial user: %w", err)
		}
	}
	return app, nil
}

func (app *App) Close() error {
	return app.DB.Close()
}

func (app *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", healthHandler)
	mux.HandleFunc("POST /api/login", app.loginHandler)
	mux.HandleFunc("POST /api/logout", app.logoutHandler)
	mux.Handle("GET /api/me", app.authenticate(http.HandlerFunc(app.meHandler)))
	mux.Handle("GET /api/users", app.authenticate(http.HandlerFunc(app.listUsersHandler)))
	mux.Handle("POST /api/users", app.authenticate(http.HandlerFunc(app.createUserHandler)))
	mux.Handle("PUT /api/users/{id}", app.authenticate(http.HandlerFunc(app.updateUserHandler)))
	mux.Handle("DELETE /api/users/{id}", app.authenticate(http.HandlerFunc(app.deleteUserHandler)))

	assets, err := fs.Sub(webFiles, "web")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, err := webFiles.ReadFile("web/index.html")
		if err != nil {
			http.Error(w, "frontend unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
