package items

import (
	"fmt"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/userItems"
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/pages/account"
)

func CardShop(
	itemsService *userItems.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		items, err := itemsService.GetItems(r)

		if err != nil {
			http.Error(
				w,
				"Failed to get user items",
				http.StatusInternalServerError,
			)
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
}

func BuyCard(
	itemsService *userItems.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(
				w,
				"Method Not Allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		cardName := r.FormValue("card")

		if cardName == "" {
			http.Error(
				w,
				"Card is required",
				http.StatusBadRequest,
			)
			return
		}

		err := itemsService.BuyCard(
			r,
			cardName,
		)

		if err != nil {
			switch err.Error() {
			case "insufficient balance":
				http.Error(
					w,
					"Недостаточно денег",
					http.StatusPaymentRequired,
				)

			case "card already owned":
				http.Error(
					w,
					"Карта уже куплена",
					http.StatusConflict,
				)

			case "card not found":
				http.Error(
					w,
					"Карта не найдена",
					http.StatusNotFound,
				)

			default:
				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
				fmt.Println(err)
			}

			return
		}

		http.Redirect(
			w,
			r,
			"/account/cards",
			http.StatusSeeOther,
		)
	}
}

func SelectCard(itemsService *userItems.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(
				w,
				"Method Not Allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		cardName := r.FormValue("card")

		if cardName == "" {
			http.Error(
				w,
				"Card is required",
				http.StatusBadRequest,
			)
			return
		}

		err := itemsService.SetCurrentCard(
			r,
			cardName,
		)

		if err != nil {
			switch err.Error() {
			case "card not owned":
				http.Error(
					w,
					"Карта не куплена",
					http.StatusForbidden,
				)

			case "card not found":
				http.Error(
					w,
					"Карта не найдена",
					http.StatusNotFound,
				)

			default:
				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}

			return
		}

		http.Redirect(
			w,
			r,
			"/account/cards",
			http.StatusSeeOther,
		)
	}
}
