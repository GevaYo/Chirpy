package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type userResponse struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type userReq struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userBody := userReq{}
	err := decoder.Decode(&userBody)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}
	user, err := cfg.db.CreateUser(r.Context(), userBody.Email)
	if err != nil {
		respondWithError(w, 500, "failed to create user", err)
	}
	userResp := userResponse{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	respondWithJSON(w, 201, userResp)
}
