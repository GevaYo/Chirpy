package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type validResponse struct {
	CleanedBody string `json:"cleaned_body"`
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
	cleaned := replaceBadWords(ch.Body)
	respondWithJSON(w, 200, validResponse{
		CleanedBody: cleaned,
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
