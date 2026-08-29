package userItems

import (
	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/repository/economics"
)

const apiURL = "http://127.0.0.1:8001/api/items"

type ItemsResponse struct {
	CurrentCard string   `json:"current_card"`
	Cards       []string `json:"cards"`
	Error       string   `json:"error"`
}

type itemsRequest struct {
	DiscordID   int    `json:"discord_id"`
	AddCard     string `json:"add_card,omitempty"`
	CurrentCard string `json:"current_card,omitempty"`
}

type Service struct {
	economics *economics.Repository
}

func NewService(
	economicsRepository *economics.Repository,
) *Service {
	return &Service{
		economics: economicsRepository,
	}
}
