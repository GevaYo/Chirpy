package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GevaYo/Chirpy/internal/auth"
	"github.com/GevaYo/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, r *http.Request) {
	records, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't fetch chirps", err)
	}

	var chirps []chirpResponse
	for _, r := range records {
		chirps = append(chirps, chirpResponse{
			Id:        r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Body:      r.Body,
			UserId:    r.UserID,
		})
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	parsedChirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 500, "failed to parse chirp ID", err)
	}

	record, err := cfg.db.GetChirpByID(r.Context(), parsedChirpId)
	if err != nil {
		respondWithError(w, 500, "couldn't fetch chirps", err)
	}

	respondWithJSON(w, 200,
		chirpResponse{
			Id:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			Body:      record.Body,
			UserId:    record.UserID,
		})
}

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type params struct {
		Body string `json:"body"`
	}
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	ch := params{}
	err = decoder.Decode(&ch) // Decode the request body
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(ch.Body) > 140 {
		respondWithError(w, 400, "chirp is too long", nil)
		return
	}
	cleaned := replaceBadWords(ch.Body)
	chirp, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams{Body: cleaned, UserID: userId})
	if err != nil {
		respondWithError(w, 500, "failed to create chirp", err)
		return
	}
	respondWithJSON(w, 201, chirpResponse{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func replaceBadWords(body string) string {
	words := strings.Split(body, " ")
	for i, w := range words {
		word := strings.ToLower(w)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
