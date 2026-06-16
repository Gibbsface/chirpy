package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type requestJSON struct {
	Body string `json:"body"`
}

type responseJSON struct {
	Valid bool `json:"valid"`
}

type errJSON struct {
	Error string `json:"error"`
}

func ApiValidateChirp(w http.ResponseWriter, r *http.Request) {

	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := requestJSON{}
	err := decoder.Decode(&reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	// at this point, we know reqJSON successfully represents the request
	if len(reqJSON.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// handle valid response
	respondWithJSON(w, http.StatusOK, responseJSON{
		Valid: true,
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, errJSON{
		Error: message,
	})
}

func ApiHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
