package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type requestJSON struct {
	Body string `json:"body"`
}

type createUserReqJSON struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type responseJSON struct {
	Cleaned_body string `json:"cleaned_body"`
}

type errJSON struct {
	Error string `json:"error"`
}

func (cfg *Config) ApiCreateUser(w http.ResponseWriter, r *http.Request) {
	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := &createUserReqJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	// at this point, we know that reqJSON is valid
	user, err := cfg.db.CreateUser(r.Context(), reqJSON.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	JSONuser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	// at this point, we know the user was created. Print the results
	fmt.Printf("User created with email %v\n", user.Email)

	//reply with JSON
	respondWithJSON(w, http.StatusCreated, JSONuser)
}

func ApiValidateChirp(w http.ResponseWriter, r *http.Request) {

	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := requestJSON{}
	err := decoder.Decode(&reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request.")
		return
	}

	// at this point, we know reqJSON successfully represents the request
	if len(reqJSON.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// at this point, we know the chirp has a valid length, now we need to clean it
	cleanChirp := cleanChirp(reqJSON.Body)

	// handle valid response
	respondWithJSON(w, http.StatusOK, responseJSON{
		Cleaned_body: cleanChirp,
	})

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
