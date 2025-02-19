package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GevaYo/Chirpy/internal/auth"
	"github.com/GevaYo/Chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type userReq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	loginBody := userReq{}
	err := decoder.Decode(&loginBody)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), loginBody.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
	}
	err = auth.CheckPasswordHash(loginBody.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
	}
	expiresIn := time.Hour

	JWT, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		respondWithError(w, 500, "couldn't create token", err)
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "couldn't create refresh token", err)
		return
	}
	err = cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{Token: refreshToken, UserID: user.ID, ExpiresAt: time.Now().Add(60 * 24 * time.Hour)})
	if err != nil {
		respondWithError(w, 500, "failed to store refresh token", err)
	}

	userResp := response{
		User: User{
			Id:         user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:      user.Email,
		},
		Token:        JWT,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, 200, userResp)
}
