package web

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, response map[string]interface{}) {
	r, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(r)
}
