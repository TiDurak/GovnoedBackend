package auth

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/config"
)

func Login(cfg config.Config, discordCfg config.DiscordConfig) http.HandlerFunc {
	oauthConfig := &oauth2.Config{
		ClientID:     discordCfg.ClientID,
		ClientSecret: discordCfg.ClientSecret,
		RedirectURL:  discordCfg.RedirectURL,

		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/v10/oauth2/token",
		},

		Scopes: []string{
			"identify",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		state, err := GenerateState()
		if err != nil {
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
			return
		}

		SetStateCookie(w, state)

		http.Redirect(
			w,
			r,
			oauthConfig.AuthCodeURL(state),
			http.StatusFound,
		)
	}
}
