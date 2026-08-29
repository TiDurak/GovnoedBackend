package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/config"
	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/database"
	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/handler"
	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/repository"
	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/service"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatalf(
			"invalid configuration: %v",
			err,
		)
	}

	projectRoot, err := findProjectRoot()

	if err != nil {
		log.Fatal(err)
	}

	if !filepath.IsAbs(cfg.DatabasePath) {
		cfg.DatabasePath = filepath.Join(projectRoot, cfg.DatabasePath)
	}

	db, err := database.Open(cfg.DatabasePath)

	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	userItemsRepository := repository.NewUserItemsRepository(db)
	userItemsService := service.NewUserItemsService(userItemsRepository)
	healthHandler := handler.NewHealthHandler(db)
	userItemsHandler := handler.NewUserItemsHandler(userItemsService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", healthHandler.Handle)
	mux.HandleFunc("/api/items", userItemsHandler.Generate)

	address := cfg.HTTPAddress + ":" + cfg.HTTPPort

	log.Printf("API listening on %s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf(
			"HTTP server stopped: %v",
			err,
		)
	}
}

func findProjectRoot() (string, error) {
	directory, err := os.Getwd()

	if err != nil {
		return "",
			fmt.Errorf(
				"failed to determine project root: %w",
				err,
			)
	}

	for {
		info, err := os.Stat(
			filepath.Join(
				directory,
				"go.mod",
			),
		)

		if err == nil && !info.IsDir() {
			return directory, nil
		}

		parent := filepath.Dir(directory)

		if parent == directory {
			return "",
				fmt.Errorf(
					"failed to locate project root from %q",
					directory,
				)
		}

		directory = parent
	}
}
