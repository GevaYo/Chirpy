package main

import (
	"encoding/json"
	"net/http"
)

type validResponse struct {
	Valid bool `json:"valid"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	err := decoder.Decode(&ch) // Decode the request body
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}

	if len(ch.Body) > 140 { // Check if the length exceeds
		respondWithError(w, 400, "chirp is too long", nil)
		return
	}

	respondWithJSON(w, 200, validResponse{
		Valid: true,
	})
}
