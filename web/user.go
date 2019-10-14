package web

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

func login(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	json.NewDecoder(r.Body).Decode(&requestData)

	if requestData["login"] == viper.GetString("login") && requestData["password"] == viper.GetString("password") {
		session, _ := store.Get(r, "user")
		session.Values["authenticated"] = true
		session.Save(r, w)
		jsonResponse(w, map[string]interface{}{"ok": true})
	} else {
		jsonResponse(w, map[string]interface{}{"ok": false})
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")
	session.Values["authenticated"] = false
	session.Save(r, w)
	jsonResponse(w, map[string]interface{}{"ok": true})
}
