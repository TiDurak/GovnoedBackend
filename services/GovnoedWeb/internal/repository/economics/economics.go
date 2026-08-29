package economics

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetBalance(userID int64) (float64, error) {
	var balance float64

	err := r.db.QueryRow(`
		SELECT balance
		FROM economics
		WHERE user_id = ?
	`, userID).Scan(&balance)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf(
			"get balance: %w",
			err,
		)
	}

	return balance, nil
}

func (r *Repository) Spend(
	userID int64,
	price int64,
) error {
	result, err := r.db.Exec(`
		UPDATE economics
		SET balance = balance - ?
		WHERE user_id = ?
		  AND balance >= ?
	`,
		price,
		userID,
		price,
	)

	if err != nil {
		return fmt.Errorf(
			"spend balance: %w",
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check spend result: %w",
			err,
		)
	}

	if rows == 0 {
		return fmt.Errorf(
			"insufficient balance",
		)
	}

	return nil
}
