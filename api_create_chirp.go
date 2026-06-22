package main

import (
	"encoding/json"
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type createChirpRequestJSON struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *Config) ApiCreateChirp(w http.ResponseWriter, r *http.Request) {
	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := &createChirpRequestJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	//validate chirp
	validatedChirp, err := validateChirp(reqJSON.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "")
	}

	// at this point, we know that reqJSON is valid
	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   validatedChirp,
		UserID: reqJSON.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}

	resJSON := chirpJSON{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	// at this point, we know the user was created. Print the results
	// fmt.Printf("User created with email %v\n", user.Email)

	//reply with JSON
	respondWithJSON(w, http.StatusCreated, resJSON)
}
