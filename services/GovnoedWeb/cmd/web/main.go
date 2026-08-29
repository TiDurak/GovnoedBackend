package main

import (
	"log"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/config"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/database"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/handler"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/handler/getPromo"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/handler/items"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/repository/economics"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/userItems"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/pages"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/pages/account"
)

type APIResponse struct {
	Key       string
	Remaining int64
	Error     string
}

func NewAPIResponse(key string, remaining int64, err string) APIResponse {
	return APIResponse{
		Key:       key,
		Remaining: remaining,
		Error:     err,
	}
}

func main() {

	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	discordCfg := config.NewDiscordConfig()
	if err := discordCfg.Validate(); err != nil {
		log.Fatalf("DiscordConfig error: %v", err)
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	economicsDb := database.Open()
	economicsRepository := economics.NewRepository(economicsDb)
	userItemsService := userItems.NewService(economicsRepository)
	mux.HandleFunc("/account/cards", items.CardShop(userItemsService))
	mux.HandleFunc("/account/card/buy", items.BuyCard(userItemsService))
	mux.HandleFunc("/account/card/select", items.SelectCard(userItemsService))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handler.Render(pages.Index()))
	mux.HandleFunc("/terms", handler.Render(pages.Terms()))
	mux.HandleFunc("/privacy", handler.Render(pages.Privacy()))
	mux.HandleFunc("/account", handler.Account)
	mux.HandleFunc("/login", auth.Login(cfg, discordCfg))
	mux.HandleFunc("/oauth/callback", auth.Callback(discordCfg))
	mux.HandleFunc("/card_shop", handler.CardShop)
	mux.HandleFunc("/generate_key", func(w http.ResponseWriter, r *http.Request) {
		response, err := getPromo.GetKey(r)
		if err != nil {
			http.Error(w, "Failed to get key", http.StatusInternalServerError)
			log.Printf("Error generating key: %v", err)
			return
		}

		err = account.GenerateKey(response).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error rendering key: %v", err)
		}
	})
	mux.HandleFunc("/logout", auth.Logout)

	address := ":" + cfg.HTTPPort

	log.Printf("GovnoedWeb listening on %s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("HTTP server stopped: %v", err)
	}
}
