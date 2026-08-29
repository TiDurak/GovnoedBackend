package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Open() *sql.DB {
	db, err := sql.Open(
		"sqlite",
		"C:/Users/ivanp/Documents/debilbot/economics.db",
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS economics (
			user_id INTEGER PRIMARY KEY,
			balance INTEGER NOT NULL DEFAULT 0,
			last_daily INTEGER
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
