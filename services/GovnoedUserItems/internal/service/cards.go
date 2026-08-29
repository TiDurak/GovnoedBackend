package service

import (
	"github.com/tidurak/GovnoedBackend/services/GovnoedUserItems/internal/repository"
)

type UserItemsService struct {
	repository *repository.UserItemsRepository
}

func NewUserItemsService(repository *repository.UserItemsRepository) *UserItemsService {
	return &UserItemsService{repository: repository}
}

func (s *UserItemsService) GetCurrentCard(discordID int64) (string, error) {
	return s.repository.GetCurrentCard(discordID)
}

func (s *UserItemsService) GetCards(discordID int64) ([]string, error) {
	return s.repository.GetCards(discordID)
}

func (s *UserItemsService) AddCard(discordID int64, card string) error {
	return s.repository.AddCard(discordID, card)
}

func (s *UserItemsService) SetCurrentCard(discordID int64, card string) error {
	return s.repository.SetCurrentCard(discordID, card)
}
