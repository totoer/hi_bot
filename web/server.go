package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user")
		value, ok := session.Values["authenticated"]
		if ok && value.(bool) {
			next.ServeHTTP(w, r)
		} else {
			response, _ := json.Marshal(map[string]interface{}{"ok": false})
			w.Write(response)
		}
	})
}

func Run() {
	http.Handle("/login", http.HandlerFunc(login))
	http.Handle("/logout", authMiddleware(http.HandlerFunc(logout)))
	http.Handle("/bot", authMiddleware(http.HandlerFunc(botHandler)))

	http.ListenAndServe(":3000", nil)
}
