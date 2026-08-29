package repository

import (
	"database/sql"
	"encoding/json"
)

type UserItemsRepository struct {
	db *sql.DB
}

func NewUserItemsRepository(db *sql.DB) *UserItemsRepository {
	return &UserItemsRepository{
		db: db,
	}
}

func (r *UserItemsRepository) GetCurrentCard(discordID int64) (string, error) {
	var currentCard string

	err := r.db.QueryRow(`
		SELECT current_card
		FROM user_items
		WHERE discord_id = ?
	`, discordID).Scan(&currentCard)

	if err == sql.ErrNoRows {
		_, err = r.db.Exec(`
			INSERT INTO user_items (
				discord_id,
				current_card,
				cards
			)
			VALUES (?, 'basic', '["basic"]')
		`, discordID)

		if err != nil {
			return "", err
		}

		return "basic", nil
	}

	if err != nil {
		return "", err
	}

	return currentCard, nil
}

func (r *UserItemsRepository) GetCards(discordID int64) ([]string, error) {
	var cardsJSON string

	err := r.db.QueryRow(`
		SELECT cards
		FROM user_items
		WHERE discord_id = ?
	`, discordID).Scan(&cardsJSON)

	if err == sql.ErrNoRows {
		_, err = r.db.Exec(`
			INSERT INTO user_items (
				discord_id,
				current_card,
				cards
			)
			VALUES (?, 'basic', '["basic"]')
		`, discordID)

		if err != nil {
			return nil, err
		}

		return []string{"basic"}, nil
	}

	if err != nil {
		return nil, err
	}

	var cards []string

	err = json.Unmarshal(
		[]byte(cardsJSON),
		&cards,
	)

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *UserItemsRepository) AddCard(discordID int64, card string) error {
	cards, err := r.GetCards(discordID)

	if err != nil {
		return err
	}

	// Не добавляем карту повторно.
	for _, ownedCard := range cards {
		if ownedCard == card {
			return nil
		}
	}

	cards = append(cards, card)

	cardsJSON, err := json.Marshal(cards)

	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		UPDATE user_items
		SET cards = ?
		WHERE discord_id = ?
	`, string(cardsJSON), discordID)

	return err
}

func (r *UserItemsRepository) SetCurrentCard(discordID int64, card string) error {
	cards, err := r.GetCards(discordID)

	if err != nil {
		return err
	}

	// Проверяем, владеет ли пользователь картой.
	for _, ownedCard := range cards {
		if ownedCard == card {
			_, err = r.db.Exec(`
				UPDATE user_items
				SET current_card = ?
				WHERE discord_id = ?
			`, card, discordID)

			return err
		}
	}

	return sql.ErrNoRows
}
