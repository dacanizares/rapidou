package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type sampleGame struct {
	Title       string
	Platform    string
	ReleaseYear int
	Description string
	SteamAppID  int
}

var sampleMuseum = []sampleGame{
	{Title: "Commander Keen 4", Platform: "MS-DOS", ReleaseYear: 1991, Description: "Mi referencia personal más alta entre los plataformeros.", SteamAppID: 9180},
	{Title: "Doom", Platform: "MS-DOS", ReleaseYear: 1993, Description: "Una experiencia excelente y favorita personal.", SteamAppID: 2280},
	{Title: "Quake", Platform: "PC", ReleaseYear: 1996, Description: "Un clásico de acción que sigue siendo una experiencia favorita.", SteamAppID: 2310},
	{Title: "Half-Life", Platform: "PC", ReleaseYear: 1998, Description: "Una experiencia narrativa y de acción valorada con 10/10.", SteamAppID: 70},
	{Title: "Portal", Platform: "PC", ReleaseYear: 2007, Description: "Breve y ligero, incluso donde todavía no alcanza el pulido de su secuela.", SteamAppID: 400},
	{Title: "Machinarium", Platform: "PC", ReleaseYear: 2009, Description: "Un conjunto exquisito de música, sonidos y narrativa sin diálogos.", SteamAppID: 40700},
}

func ensureSampleMuseum(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	createdAt := time.Now().UTC().Format(time.RFC3339)

	for _, sample := range sampleMuseum {
		var gameID int64
		err := tx.QueryRow(`SELECT id FROM games WHERE title = ? AND platform = ? LIMIT 1`, sample.Title, sample.Platform).Scan(&gameID)
		if errors.Is(err, sql.ErrNoRows) {
			result, insertErr := tx.Exec(`
				INSERT INTO games (title, platform, release_year, description, created_at)
				VALUES (?, ?, ?, ?, ?)
			`, sample.Title, sample.Platform, sample.ReleaseYear, sample.Description, createdAt)
			if insertErr != nil {
				return insertErr
			}
			gameID, err = result.LastInsertId()
		}
		if err != nil {
			return err
		}

		urls := []string{
			fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%d/header.jpg", sample.SteamAppID),
			fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%d/library_hero.jpg", sample.SteamAppID),
		}
		for index, remoteURL := range urls {
			if _, err := tx.Exec(`
				INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
				VALUES (?, ?, ?, ?)
			`, gameID, remoteURL, fmt.Sprintf("%s — imagen %d", sample.Title, index+1), createdAt); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
