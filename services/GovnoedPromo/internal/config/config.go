package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPAddress  string
	HTTPPort     string
	DatabasePath string

	PromoReward int
}

func Load() Config {
	return Config{
		HTTPAddress: "127.0.0.1",
		HTTPPort:    "8000",

		DatabasePath: getEnv("PROMO_DATABASE_PATH", "./data/database.db"),
		PromoReward:  750,
	}
}

func (c Config) Validate() error {
	if c.HTTPPort == "" {
		return errors.New("HTTP port is not configured")
	}

	if c.DatabasePath == "" {
		return errors.New("database path is not configured")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
