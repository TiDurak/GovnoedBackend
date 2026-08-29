package auth

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ClearSession(w, r)

	http.Redirect(w, r, "/", http.StatusFound)
}
