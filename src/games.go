package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Game struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Platform    string      `json:"platform"`
	ReleaseYear int         `json:"release_year"`
	Description string      `json:"description"`
	SteamURL    string      `json:"steam_url"`
	GOGURL      string      `json:"gog_url"`
	CreatedAt   string      `json:"created_at"`
	Images      []GameImage `json:"images"`
}

type gameInput struct {
	Title       string `json:"title"`
	Platform    string `json:"platform"`
	ReleaseYear int    `json:"release_year"`
	Description string `json:"description"`
	SteamURL    string `json:"steam_url"`
	GOGURL      string `json:"gog_url"`
}

func (app *App) listGamesHandler(w http.ResponseWriter, _ *http.Request) {
	games, err := listGames(app.DB)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list museum pieces")
		return
	}
	writeJSON(w, http.StatusOK, games)
}

func (app *App) createGameHandler(w http.ResponseWriter, r *http.Request) {
	var input gameInput
	if !readJSON(w, r, &input) {
		return
	}
	normalizeGameInput(&input)
	if message := validateGame(input); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return
	}
	game, err := insertGame(app.DB, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create museum piece")
		return
	}
	writeJSON(w, http.StatusCreated, game)
}

func (app *App) updateGameHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid museum piece ID")
		return
	}
	var input gameInput
	if !readJSON(w, r, &input) {
		return
	}
	normalizeGameInput(&input)
	if message := validateGame(input); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return
	}
	game, err := updateGame(app.DB, id, input)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "museum piece not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update museum piece")
		return
	}
	writeJSON(w, http.StatusOK, game)
}

func (app *App) deleteGameHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid museum piece ID")
		return
	}
	result, err := app.DB.Exec(`DELETE FROM games WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete museum piece")
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		writeError(w, http.StatusNotFound, "museum piece not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func normalizeGameInput(input *gameInput) {
	input.Title = strings.TrimSpace(input.Title)
	input.Platform = strings.TrimSpace(input.Platform)
	input.Description = strings.TrimSpace(input.Description)
	input.SteamURL = strings.TrimSpace(input.SteamURL)
	input.GOGURL = strings.TrimSpace(input.GOGURL)
}

func validateGame(input gameInput) string {
	if input.Title == "" {
		return "title is required"
	}
	if input.Platform == "" {
		return "platform is required"
	}
	if input.ReleaseYear < 1950 || input.ReleaseYear > time.Now().Year()+1 {
		return "release year is invalid"
	}
	if !validStoreURL(input.SteamURL, "store.steampowered.com") {
		return "Steam URL must be an HTTPS Steam store URL"
	}
	if !validStoreURL(input.GOGURL, "gog.com") {
		return "GOG URL must be an HTTPS GOG URL"
	}
	return ""
}

func validStoreURL(value, domain string) bool {
	if value == "" {
		return true
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return host == domain || strings.HasSuffix(host, "."+domain)
}

func insertGame(db *sql.DB, input gameInput) (Game, error) {
	createdAt := time.Now().UTC().Format(time.RFC3339)
	result, err := db.Exec(`
		INSERT INTO games (title, platform, release_year, description, steam_url, gog_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, input.Title, input.Platform, input.ReleaseYear, input.Description, input.SteamURL, input.GOGURL, createdAt)
	if err != nil {
		return Game{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Game{}, err
	}
	return Game{ID: id, Title: input.Title, Platform: input.Platform, ReleaseYear: input.ReleaseYear, Description: input.Description, SteamURL: input.SteamURL, GOGURL: input.GOGURL, CreatedAt: createdAt, Images: []GameImage{}}, nil
}

func updateGame(db *sql.DB, id int64, input gameInput) (Game, error) {
	result, err := db.Exec(`
		UPDATE games SET title = ?, platform = ?, release_year = ?, description = ?, steam_url = ?, gog_url = ? WHERE id = ?
	`, input.Title, input.Platform, input.ReleaseYear, input.Description, input.SteamURL, input.GOGURL, id)
	if err != nil {
		return Game{}, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return Game{}, sql.ErrNoRows
	}
	return gameByID(db, id)
}

func gameByID(db *sql.DB, id int64) (Game, error) {
	var game Game
	err := db.QueryRow(`
		SELECT id, title, platform, release_year, description, steam_url, gog_url, created_at FROM games WHERE id = ?
	`, id).Scan(&game.ID, &game.Title, &game.Platform, &game.ReleaseYear, &game.Description, &game.SteamURL, &game.GOGURL, &game.CreatedAt)
	if err != nil {
		return Game{}, err
	}
	game.Images, err = imagesByGameID(db, game.ID)
	return game, err
}

func listGames(db *sql.DB) ([]Game, error) {
	rows, err := db.Query(`
		SELECT id, title, platform, release_year, description, steam_url, gog_url, created_at
		FROM games ORDER BY release_year, title
	`)
	if err != nil {
		return nil, err
	}
	games := []Game{}
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Platform, &game.ReleaseYear, &game.Description, &game.SteamURL, &game.GOGURL, &game.CreatedAt); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for index := range games {
		games[index].Images, err = imagesByGameID(db, games[index].ID)
		if err != nil {
			return nil, err
		}
	}
	return games, nil
}
