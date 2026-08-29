package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/config"
	"github.com/tidurak/GovnoedBackend/GovnoedPromoAPI/internal/repository"
)

const (
	keyLength   = 12
	keyAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type PromoService struct {
	repository *repository.PromoRepository
}

func NewPromoService(
	repository *repository.PromoRepository,
) *PromoService {

	return &PromoService{
		repository: repository,
	}
}

func (s *PromoService) GenerateKey(
	discordID int64,
) (string, error) {
	cfg := config.Load()
	// no cfg.Validate() check here because it was already validated in main.go
	rewardAmount := cfg.PromoReward
	key, err := generateKey()

	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(
		[]byte(key),
	)

	keyHash := hex.EncodeToString(
		hash[:],
	)

	err = s.repository.SaveKey(
		keyHash,
		discordID,
		rewardAmount,
	)

	if err != nil {
		return "", err
	}

	return key, nil
}

func generateKey() (string, error) {

	result := make(
		[]byte,
		keyLength,
	)

	for i := range result {

		index, err := rand.Int(
			rand.Reader,
			big.NewInt(
				int64(len(keyAlphabet)),
			),
		)

		if err != nil {
			return "", err
		}

		result[i] =
			keyAlphabet[index.Int64()]
	}

	return string(result), nil
}
