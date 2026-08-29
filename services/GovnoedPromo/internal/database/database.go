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

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
