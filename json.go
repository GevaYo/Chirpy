package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	respondWithJSON(w, statusCode, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
