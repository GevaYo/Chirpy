package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/GevaYo/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	refreshTokenRecord, err := cfg.db.GetToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 500, "failed to fetch token record", err)
		return
	}
	if refreshTokenRecord.ExpiresAt.Before(time.Now()) || refreshTokenRecord.RevokedAt.Valid {
		respondWithError(w, 401, "Unauthorized", errors.New("refresh token expired or was revoked"))
		return
	}
	jwt, err := auth.MakeJWT(refreshTokenRecord.UserID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "failed to create access token", err)
		return
	}

	respondWithJSON(w, 200, response{Token: jwt})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	err = cfg.db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 500, "failed to revoke token", err)
		return
	}
	respondWithJSON(w, 204, nil)
}
