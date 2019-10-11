package web

import (
	"encoding/json"

	"hi_bot/models"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func get(w http.ResponseWriter, r *http.Request) {
	dbCLient := r.Context().Value("dbClient").(*mongo.Client)
	bots := models.FindAllBot(dbCLient)

	result, err := json.Marshal(bots)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func post(w http.ResponseWriter, r *http.Request) {}

func delete(w http.ResponseWriter, r *http.Request) {}
