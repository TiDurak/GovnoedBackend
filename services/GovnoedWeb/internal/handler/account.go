package handler

import (
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/pages/account"
)

func Account(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetSession(r)

	if !ok {
		if err := account.Guest().Render(r.Context(), w); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := account.Account(
		user.Username,
		user.DiscordID,
	).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
