package main

import (
	"encoding/json"
	"net/http"

	"github.com/GevaYo/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type request struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) webhooksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		respondWithError(w, 401, "Unauthorized", err)
	}

	decoder := json.NewDecoder(r.Body)
	body := request{}
	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if body.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userId, _ := uuid.Parse(body.Data.UserID)
	_, err = cfg.db.GetUserById(r.Context(), userId)
	if err != nil {
		respondWithError(w, 404, "user ID doesn't exist", err)
	}
	err = cfg.db.MakeChirpyRed(r.Context(), userId)
	if err != nil {
		respondWithError(w, 400, "invalid user ID", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
