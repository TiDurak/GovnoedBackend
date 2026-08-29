package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/config"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/repository"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/service"
)

type PromoHandler struct {
	service *service.PromoService
}

func NewPromoHandler(
	service *service.PromoService,
) *PromoHandler {
	return &PromoHandler{
		service: service,
	}
}

func (h *PromoHandler) Generate(
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
				"error": "invalid_discord_id_or_json",
			},
		)

		return
	}

	key, err := h.service.GenerateKey(
		request.DiscordID,
	)

	if err != nil {
		var cooldownErr *repository.CooldownError
		if errors.As(err, &cooldownErr) {

			writeJSON(
				w,
				http.StatusTooManyRequests,
				map[string]any{
					"error":     "cooldown",
					"remaining": cooldownErr.Remaining,
				},
			)

			return
		}

		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal_server_error",
			},
		)

		return
	}

	cfg := config.Load()
	// no cfg.Validate() check here because it was already validated in main.go

	writeJSON(
		w,
		http.StatusOK,
		generateResponse{
			Key:    key,
			Reward: cfg.PromoReward,
		},
	)
}
