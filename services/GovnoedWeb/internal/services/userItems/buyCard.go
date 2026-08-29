package userItems

import (
	"fmt"
	"net/http"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/services/auth"
)

func (s *Service) BuyCard(
	r *http.Request,
	cardName string,
) error {
	user, ok := auth.GetSession(r)
	if !ok {
		return fmt.Errorf("user is not authenticated")
	}

	card := findCard(cardName)

	if card == nil {
		return fmt.Errorf("card not found")
	}

	// Проверяем, нет ли карты уже у пользователя.
	items, err := s.GetItems(r)
	if err != nil {
		return fmt.Errorf(
			"get user items: %w",
			err,
		)
	}

	for _, ownedCard := range items.Cards {
		if ownedCard == card.Name {
			return fmt.Errorf("card already owned")
		}
	}

	// Бесплатная карта.
	if card.Price == 0 {
		return s.addCard(
			r,
			user.DiscordID,
			card.Name,
		)
	}

	// Проверяем баланс через существующий repository.
	balance, err := s.economics.GetBalance(
		int64(user.DiscordID),
	)
	if err != nil {
		return fmt.Errorf(
			"get balance: %w",
			err,
		)
	}

	if balance < float64(card.Price) {
		return fmt.Errorf("insufficient balance")
	}

	// Списываем деньги через economics repository.
	err = s.economics.Spend(
		int64(user.DiscordID),
		int64(card.Price),
	)
	if err != nil {
		return fmt.Errorf(
			"remove balance: %w",
			err,
		)
	}

	// Выдаём карту.
	err = s.addCard(
		r,
		user.DiscordID,
		card.Name,
	)
	if err != nil {
		return fmt.Errorf(
			"add card: %w",
			err,
		)
	}

	return nil
}
