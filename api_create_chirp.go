package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
	"github.com/Gibbsface/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type createChirpRequestJSON struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (c *Config) ApiCreateChirp(w http.ResponseWriter, r *http.Request) {
	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := &createChirpRequestJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	//validate JWT
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	id, err := auth.ValidateJWT(jwt, c.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT implied the wrong uuid")
		return
	}

	//validate chirp
	validatedChirp, err := validateChirp(reqJSON.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	// at this point, we know that reqJSON is valid
	dbChirp, err := c.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   validatedChirp,
		UserID: id,
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

	// at this point, we know the chirp was created. Print the results
	fmt.Printf("Chirp created by %v\n", dbChirp.ID)

	//reply with JSON
	respondWithJSON(w, http.StatusCreated, resJSON)
}
