package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
)

const sessionCookieName = "session_id"

type Session struct {
	DiscordID int
	Username  string
}

var (
	sessions   = make(map[string]Session)
	sessionsMu sync.RWMutex
)

func CreateSession(w http.ResponseWriter, user User) error {
	sessionIDBytes := make([]byte, 32)

	if _, err := rand.Read(sessionIDBytes); err != nil {
		return fmt.Errorf("generate session ID: %w", err)
	}

	sessionID := base64.RawURLEncoding.EncodeToString(sessionIDBytes)

	sessionsMu.Lock()
	sessions[sessionID] = Session{
		DiscordID: user.DiscordID,
		Username:  user.Username,
	}
	sessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30,
	})

	return nil
}

func GetSession(r *http.Request) (User, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return User{}, false
	}

	sessionsMu.RLock()
	session, ok := sessions[cookie.Value]
	sessionsMu.RUnlock()

	if !ok {
		return User{}, false
	}

	return User{
		DiscordID: session.DiscordID,
		Username:  session.Username,
	}, true
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		sessionsMu.Lock()
		delete(sessions, cookie.Value)
		sessionsMu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
