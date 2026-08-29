package userItems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
)

func (s *Service) SetCurrentCard(
	r *http.Request,
	cardName string,
) error {
	user, ok := auth.GetSession(r)
	if !ok {
		return fmt.Errorf("user is not authenticated")
	}

	// Проверяем, существует ли такая карта вообще.
	if findCard(cardName) == nil {
		return fmt.Errorf("card not found")
	}

	// Проверяем, есть ли карта у пользователя.
	items, err := s.GetItems(r)
	if err != nil {
		return fmt.Errorf(
			"get user items: %w",
			err,
		)
	}

	owned := false

	for _, ownedCard := range items.Cards {
		if ownedCard == cardName {
			owned = true
			break
		}
	}

	if !owned {
		return fmt.Errorf("card not owned")
	}

	// Отправляем запрос в UserItems API
	// для изменения выбранной карты.
	payload, err := json.Marshal(itemsRequest{
		DiscordID:   user.DiscordID,
		CurrentCard: cardName,
	})
	if err != nil {
		return fmt.Errorf(
			"marshal current card request: %w",
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
		return fmt.Errorf(
			"create current card request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"current card API request: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result ItemsResponse

		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			if result.Error != "" {
				return fmt.Errorf(
					"current card API: %s",
					result.Error,
				)
			}
		}

		return fmt.Errorf(
			"current card API returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
