package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *Config) ApiGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error: could not parse chirp ID")
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error: while fetching chirp.")
	} else {
		respondWithJSON(w, 200, chirpJSON{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
}
