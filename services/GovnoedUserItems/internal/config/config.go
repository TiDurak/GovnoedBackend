package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPAddress  string
	HTTPPort     string
	DatabasePath string
}

func Load() Config {
	return Config{
		HTTPAddress: "127.0.0.1",
		HTTPPort:    getEnv("GUITEMS_HTTP_PORT", "8001"),

		DatabasePath: getEnv("CARDS_DATABASE_PATH", "C:/Users/ivanp/Documents/debilbot/economics.db"),
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
