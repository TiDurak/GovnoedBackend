package auth

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"

	"github.com/tidurak/GovnoedBackend/services/GovnoedWeb/internal/config"
)

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type User struct {
	DiscordID int
	Username  string
}

func NewUser(discordID int, username string) User {
	return User{
		DiscordID: discordID,
		Username:  username,
	}
}

func Callback(discordCfg config.DiscordConfig) http.HandlerFunc {
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
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if code == "" || state == "" {
			http.Error(
				w,
				"Ошибка авторизации.",
				http.StatusBadRequest,
			)
			return
		}

		savedState, err := GetStateCookie(r)
		if err != nil {
			http.Error(
				w,
				"Недействительный OAuth state.",
				http.StatusForbidden,
			)
			return
		}

		if state != savedState {
			http.Error(
				w,
				"Недействительный OAuth state.",
				http.StatusForbidden,
			)
			return
		}

		ClearStateCookie(w)

		token, err := oauthConfig.Exchange(
			r.Context(),
			code,
		)
		if err != nil {
			http.Error(
				w,
				"Не удалось получить access token.",
				http.StatusInternalServerError,
			)
			return
		}

		client := oauthConfig.Client(
			r.Context(),
			token,
		)

		resp, err := client.Get(
			"https://discord.com/api/v10/users/@me",
		)
		if err != nil {
			http.Error(
				w,
				"Не удалось получить данные Discord.",
				http.StatusInternalServerError,
			)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(
				w,
				"Discord вернул ошибку.",
				http.StatusInternalServerError,
			)
			return
		}

		var user DiscordUser

		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			http.Error(
				w,
				"Не удалось обработать данные Discord.",
				http.StatusInternalServerError,
			)
			return
		}

		if user.ID == "" {
			http.Error(
				w,
				"Не удалось получить Discord ID.",
				http.StatusInternalServerError,
			)
			return
		}
		discordID, err := strconv.Atoi(user.ID)
		if err != nil {
			http.Error(
				w,
				"Не удалось преобразовать Discord ID.",
				http.StatusInternalServerError,
			)
			return
		}
		authUser := NewUser(
			discordID,
			user.Username,
		)

		if err := CreateSession(w, authUser); err != nil {
			http.Error(
				w,
				"Не удалось создать сессию.",
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(
			w,
			r,
			"/account",
			http.StatusFound,
		)
	}
}
