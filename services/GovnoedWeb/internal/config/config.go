package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPAddress          string
	HTTPPort             string
	GovnoedPromoHTTPPort string
	EconomicsDbPath      string
}

func NewConfig() Config {
	return Config{
		HTTPAddress:          getEnv("GWEB_HTTP_ADDRESS", "127.0.0.1"),
		HTTPPort:             getEnv("GWEB_HTTP_PORT", "8080"),
		GovnoedPromoHTTPPort: getEnv("GPROMO_HTTP_PORT", "8000"),
		EconomicsDbPath:      getEnv("ECONOMICS_DATABASE_PATH", ""),
	}
}

func (c Config) Validate() error {
	if c.HTTPAddress == "" {
		return errors.New("HTTP address is not configured")
	}
	if c.HTTPPort == "" {
		return errors.New("HTTP port is not configured")
	}

	return nil
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func NewDiscordConfig() DiscordConfig {
	return DiscordConfig{
		ClientID:     getEnv("DISCORD_CLIENT_ID", ""),
		ClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
		RedirectURL:  getEnv("DISCORD_REDIRECT_URL", ""),
	}
}

func (c DiscordConfig) Validate() error {
	if c.ClientID == "" {
		return errors.New("Discord client ID is not configured")
	}
	if c.ClientSecret == "" {
		return errors.New("Discord client secret is not configured")
	}
	if c.RedirectURL == "" {
		return errors.New("Discord redirect URL is not configured")
	}
	return nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists || value == "" {
		return fallback
	}

	return value
}
