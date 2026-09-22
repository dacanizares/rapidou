package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const maxImageSize = 5 << 20

type GameImage struct {
	ID  int64  `json:"id"`
	Src string `json:"src"`
	Alt string `json:"alt"`
}

func (app *App) uploadGameImageHandler(w http.ResponseWriter, r *http.Request) {
	gameID, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid museum piece ID")
		return
	}
	game, err := gameByID(app.DB, gameID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "museum piece not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read museum piece")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize+(1<<20))
	if err := r.ParseMultipartForm(maxImageSize + (1 << 20)); err != nil {
		writeError(w, http.StatusBadRequest, "image upload is too large or invalid")
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "image is required")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxImageSize+1))
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read image")
		return
	}
	if len(data) == 0 || len(data) > maxImageSize {
		writeError(w, http.StatusBadRequest, "image must be 5 MB or smaller")
		return
	}
	mimeType := http.DetectContentType(data)
	if !allowedImageType(mimeType) {
		writeError(w, http.StatusBadRequest, "image must be JPEG, PNG, GIF, or WebP")
		return
	}
	alt := strings.TrimSpace(r.FormValue("alt"))
	if alt == "" {
		alt = game.Title
	}
	result, err := app.DB.Exec(`
		INSERT INTO game_images (game_id, mime_type, image_data, alt_text, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, gameID, mimeType, data, alt, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save image")
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save image")
		return
	}
	writeJSON(w, http.StatusCreated, GameImage{ID: id, Src: fmt.Sprintf("/api/game-images/%d", id), Alt: alt})
}

func (app *App) gameImageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	var mimeType string
	var data []byte
	err := app.DB.QueryRow(`
		SELECT mime_type, image_data FROM game_images WHERE id = ? AND image_data IS NOT NULL
	`, id).Scan(&mimeType, &data)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "image unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(data)
}

func (app *App) deleteGameImageHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid image ID")
		return
	}
	result, err := app.DB.Exec(`DELETE FROM game_images WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete image")
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func imagesByGameID(db *sql.DB, gameID int64) ([]GameImage, error) {
	rows, err := db.Query(`SELECT id, url, alt_text FROM game_images WHERE game_id = ? ORDER BY id`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	images := []GameImage{}
	for rows.Next() {
		var image GameImage
		var remoteURL sql.NullString
		if err := rows.Scan(&image.ID, &remoteURL, &image.Alt); err != nil {
			return nil, err
		}
		if remoteURL.Valid {
			image.Src = remoteURL.String
		} else {
			image.Src = fmt.Sprintf("/api/game-images/%d", image.ID)
		}
		images = append(images, image)
	}
	return images, rows.Err()
}

func allowedImageType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif" || mimeType == "image/webp"
}
