package main

import (
	"encoding/json"
	"net/http"

	"github.com/GevaYo/Chirpy/internal/auth"
	"github.com/GevaYo/Chirpy/internal/database"
)

func (cfg *apiConfig) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type userReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorzied", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	body := userReq{}
	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorzied", err)
		return
	}

	hashedPwd, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{HashedPassword: hashedPwd, Email: body.Email, ID: userId})
	if err != nil {
		respondWithError(w, 500, "failed to update user", err)
	}

	respondWithJSON(w, 200, response{
		User: User{
			Id:          updatedUser.ID,
			Created_at:  updatedUser.CreatedAt,
			Updated_at:  updatedUser.UpdatedAt,
			Email:       updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
		},
	})
}
