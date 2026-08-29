package userItems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
)

func (s *Service) GetItems(r *http.Request) (ItemsResponse, error) {
	user, ok := auth.GetSession(r)
	if !ok {
		return ItemsResponse{}, fmt.Errorf(
			"user is not authenticated",
		)
	}

	payload, err := json.Marshal(itemsRequest{
		DiscordID: user.DiscordID,
	})
	if err != nil {
		return ItemsResponse{}, fmt.Errorf(
			"marshal request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		apiURL,
		bytes.NewReader(payload),
	)
	if err != nil {
		return ItemsResponse{}, fmt.Errorf(
			"create request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ItemsResponse{}, fmt.Errorf(
			"items API request: %w",
			err,
		)
	}
	defer resp.Body.Close()

	var result ItemsResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ItemsResponse{}, fmt.Errorf(
			"decode items response (status %d): %w",
			resp.StatusCode,
			err,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf(
			"items API returned status %d: %s",
			resp.StatusCode,
			result.Error,
		)
	}

	return result, nil
}
