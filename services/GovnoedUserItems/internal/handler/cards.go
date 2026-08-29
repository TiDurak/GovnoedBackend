package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/service"
)

type UserItemsHandler struct {
	service *service.UserItemsService
}

func NewUserItemsHandler(
	service *service.UserItemsService,
) *UserItemsHandler {
	return &UserItemsHandler{
		service: service,
	}
}

func (h *UserItemsHandler) Generate(
	w http.ResponseWriter,
	r *http.Request,
) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "failed_to_read_body",
			},
		)

		return
	}

	var request generateRequest

	err = json.Unmarshal(body, &request)

	if err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": fmt.Sprintf(
					"invalid_json: %v",
					err,
				),
			},
		)

		return
	}

	if request.DiscordID <= 0 {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid_discord_id",
			},
		)

		return
	}

	if request.AddCard != nil {
		if *request.AddCard == "" {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": "invalid_card",
				},
			)

			return
		}

		err = h.service.AddCard(
			request.DiscordID,
			*request.AddCard,
		)

		if err != nil {
			writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{
					"error": "internal_server_error",
				},
			)
			fmt.Print(err)

			return
		}
	}

	if request.CurrentCard != nil {
		if *request.CurrentCard == "" {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": "invalid_current_card",
				},
			)

			return
		}

		err = h.service.SetCurrentCard(
			request.DiscordID,
			*request.CurrentCard,
		)

		if err != nil {
			fmt.Println(err)

			switch err.Error() {
			case "card not owned":
				writeJSON(
					w,
					http.StatusForbidden,
					map[string]string{
						"error": "card_not_owned",
					},
				)

			case "card not found":
				writeJSON(
					w,
					http.StatusNotFound,
					map[string]string{
						"error": "card_not_found",
					},
				)

			default:
				writeJSON(
					w,
					http.StatusInternalServerError,
					map[string]string{
						"error": "internal_server_error",
					},
				)
			}

			return
		}
	}

	currentCard, err := h.service.GetCurrentCard(
		request.DiscordID,
	)

	if err != nil {
		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal_server_error",
			},
		)
		fmt.Println(err)

		return
	}

	cards, err := h.service.GetCards(
		request.DiscordID,
	)

	if err != nil {
		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal_server_error",
			},
		)
		fmt.Println(err)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		generateResponse{
			CurrentCard: currentCard,
			Cards:       cards,
		},
	)
}
