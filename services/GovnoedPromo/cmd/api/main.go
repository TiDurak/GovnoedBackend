package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/config"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/database"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/handler"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/repository"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/service"
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

	if err := database.Migrate(db); err != nil {
		log.Fatalf(
			"failed to migrate database: %v",
			err,
		)
	}

	promoRepository := repository.NewPromoRepository(db)
	promoService := service.NewPromoService(promoRepository)
	healthHandler := handler.NewHealthHandler(db)
	promoHandler := handler.NewPromoHandler(promoService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", healthHandler.Handle)
	mux.HandleFunc("/api/promo/generate", promoHandler.Generate)

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
