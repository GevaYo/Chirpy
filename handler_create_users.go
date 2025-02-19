package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GevaYo/Chirpy/internal/auth"
	"github.com/GevaYo/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type userReq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	userBody := userReq{}
	err := decoder.Decode(&userBody)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}
	hashedPwd, err := auth.HashPassword(userBody.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password", err)
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: userBody.Email, HashedPassword: hashedPwd})
	if err != nil {
		respondWithError(w, 500, "failed to create user", err)
	}
	userResp := response{
		User: User{
			Id:         user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:      user.Email,
		},
	}

	respondWithJSON(w, 201, userResp)
}
