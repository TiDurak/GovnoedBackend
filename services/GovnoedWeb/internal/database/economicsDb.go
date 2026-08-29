package database

import (
	"database/sql"
	"log"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/config"

	_ "modernc.org/sqlite"
)

var cfg = config.NewConfig()

func Open() *sql.DB {
	db, err := sql.Open(
		"sqlite",
		cfg.EconomicsDbPath,
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
