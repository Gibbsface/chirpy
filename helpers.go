package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type errJSON struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, errJSON{
		Error: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	// attempt to marshal
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}

	return cleanChirp(chirp), nil
}

func cleanChirp(chirp string) string {
	tokens := strings.Split(chirp, " ")

	profanity := make(map[string]struct{})
	profanity["kerfuffle"] = struct{}{}
	profanity["sharbert"] = struct{}{}
	profanity["fornax"] = struct{}{}

	for i, t := range tokens {
		if _, ok := profanity[strings.ToLower(t)]; ok {
			tokens[i] = "****"
		}
	}

	return strings.Join(tokens, " ")
}
