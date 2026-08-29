package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(db *sql.DB) error {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at INTEGER NOT NULL
		)
	`)

	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(
		migrations,
		"migrations",
	)

	if err != nil {
		return err
	}

	sort.Slice(
		entries,
		func(i, j int) bool {
			return entries[i].Name() <
				entries[j].Name()
		},
	)

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		id, err := migrationID(name)

		if err != nil {
			return err
		}

		var exists bool

		err = db.QueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM migrations
				WHERE id = ?
			)
		`, id).Scan(&exists)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		content, err := fs.ReadFile(
			migrations,
			"migrations/"+name,
		)

		if err != nil {
			return err
		}

		tx, err := db.Begin()

		if err != nil {
			return err
		}

		_, err = tx.Exec(
			string(content),
		)

		if err != nil {
			tx.Rollback()

			return fmt.Errorf(
				"migration %s failed: %w",
				name,
				err,
			)
		}

		_, err = tx.Exec(`
			INSERT INTO migrations (
				id,
				name,
				applied_at
			)
			VALUES (?, ?, unixepoch())
		`, id, name)

		if err != nil {
			tx.Rollback()

			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func migrationID(name string) (int, error) {

	parts := strings.SplitN(
		name,
		"_",
		2,
	)

	if len(parts) != 2 {
		return 0, fmt.Errorf(
			"invalid migration name: %s",
			name,
		)
	}

	id, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, fmt.Errorf(
			"invalid migration ID in %s: %w",
			name,
			err,
		)
	}

	return id, nil
}
