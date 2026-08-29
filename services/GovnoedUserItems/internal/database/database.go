package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	directory := filepath.Dir(path)

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(
		"sqlite",
		path,
	)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_items (
		discord_id INTEGER PRIMARY KEY,
		current_card TEXT NOT NULL DEFAULT 'basic',
		cards TEXT NOT NULL DEFAULT '["basic"]'
	)`)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
