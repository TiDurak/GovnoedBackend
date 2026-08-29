package userItems

import (
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/data"
)

func findCard(name string) *data.Card {
	for i := range data.Cards {
		if data.Cards[i].Name == name {
			return &data.Cards[i]
		}
	}

	return nil
}
