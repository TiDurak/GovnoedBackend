package getPromo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
)

type PromoResponse struct {
	Key       string `json:"key"`
	Reward    int    `json:"reward"`
	Error     string `json:"error"`
	Remaining int    `json:"remaining"`
}

func GetKey(r *http.Request) (PromoResponse, error) {
	user, ok := auth.GetSession(r)
	if !ok {
		return PromoResponse{}, fmt.Errorf("user is not authenticated")
	}

	payload, err := json.Marshal(struct {
		DiscordID int `json:"discord_id"`
	}{
		DiscordID: user.DiscordID,
	})
	if err != nil {
		return PromoResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		"http://127.0.0.1:8000/api/promo/generate",
		bytes.NewReader(payload),
	)
	if err != nil {
		return PromoResponse{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return PromoResponse{}, fmt.Errorf("promo API request: %w", err)
	}
	defer resp.Body.Close()

	var result PromoResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PromoResponse{}, fmt.Errorf(
			"decode promo response (status %d): %w",
			resp.StatusCode,
			err,
		)
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusBadRequest, http.StatusTooManyRequests:
		return result, nil

	default:
		return PromoResponse{}, fmt.Errorf(
			"promo API returned unexpected status %d",
			resp.StatusCode,
		)
	}
}
