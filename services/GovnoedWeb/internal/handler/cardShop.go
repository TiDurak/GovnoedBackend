package handler

import (
	"fmt"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/userItems"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/pages/account"
)

var ItemsService *userItems.Service

func CardShop(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetSession(r)

	if !ok {
		if err := account.Guest().Render(
			r.Context(),
			w,
		); err != nil {
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
			fmt.Println(err)
		}

		return
	}

	items, err := ItemsService.GetItems(r)
	if err != nil {
		http.Error(
			w,
			"Failed to get user items",
			http.StatusInternalServerError,
		)
		fmt.Println(err)
		return
	}

	if err := account.CardShop(
		user.Username,
		user.DiscordID,
		items,
	).Render(r.Context(), w); err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		fmt.Println(err)
	}
}
