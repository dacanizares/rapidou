package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE COLLATE NOCASE,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    platform TEXT NOT NULL,
    release_year INTEGER NOT NULL,
    description TEXT NOT NULL,
    steam_url TEXT NOT NULL DEFAULT '',
    gog_url TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS game_images (
    id INTEGER PRIMARY KEY,
    game_id INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    url TEXT,
    mime_type TEXT,
    image_data BLOB,
    alt_text TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE (game_id, url),
    CHECK ((url IS NOT NULL AND image_data IS NULL) OR (url IS NULL AND image_data IS NOT NULL))
);
`

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	for _, migration := range []string{
		`ALTER TABLE games ADD COLUMN steam_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE games ADD COLUMN gog_url TEXT NOT NULL DEFAULT ''`,
	} {
		if _, err := db.Exec(migration); err != nil && !isDuplicateColumnError(err) {
			db.Close()
			return nil, fmt.Errorf("migrate database: %w", err)
		}
	}
	return db, nil
}

func isDuplicateColumnError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate column name")
}
