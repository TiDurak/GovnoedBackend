package userItems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Service) addCard(r *http.Request, discordID int, cardName string) error {
	payload, err := json.Marshal(itemsRequest{
		DiscordID: discordID,
		AddCard:   cardName,
	})
	if err != nil {
		return fmt.Errorf(
			"marshal add card request: %w",
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
			"create add card request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"add card API request: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"add card API returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
