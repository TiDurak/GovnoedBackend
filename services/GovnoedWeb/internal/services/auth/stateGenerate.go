package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateState() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate OAuth state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
