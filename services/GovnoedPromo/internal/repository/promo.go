package repository

import (
	"database/sql"
	"time"
)

type CooldownError struct {
	Remaining int64
}

func (e *CooldownError) Error() string {
	return "promo key cooldown is active"
}

type PromoRepository struct {
	db *sql.DB
}

func NewPromoRepository(db *sql.DB) *PromoRepository {
	return &PromoRepository{
		db: db,
	}
}

func (r *PromoRepository) GetCreatedAt(
	creatorID int64,
) (int64, error) {

	var createdAt int64

	err := r.db.QueryRow(`
		SELECT created_at
		FROM redeem_keys
		WHERE creator_id = ?
	`, creatorID).Scan(&createdAt)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return createdAt, nil
}

func (r *PromoRepository) SaveKey(
	keyHash string,
	creatorID int64,
	reward int,
) error {

	now := time.Now().Unix()

	createdAt, err := r.GetCreatedAt(creatorID)

	if err != nil {
		return err
	}

	if createdAt != 0 {
		// 12 hours cooldown check
		if now-createdAt < 12*60*60 {
			remaining := 12*60*60 - (now - createdAt)
			return &CooldownError{Remaining: remaining}
		}
	}

	_, err = r.db.Exec(`
		INSERT INTO redeem_keys (
			key_hash,
			creator_id,
			created_at,
			is_used,
			reward
		)
		VALUES (?, ?, ?, FALSE, ?)

		ON CONFLICT(creator_id)
		DO UPDATE SET
			key_hash = excluded.key_hash,
			created_at = excluded.created_at,
			is_used = FALSE,
			reward = excluded.reward
	`,
		keyHash,
		creatorID,
		now,
		reward,
	)

	return err
}
