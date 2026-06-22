package main

import (
	"net/http"
)

func (cfg *Config) ApiGetAllChirps(w http.ResponseWriter, r *http.Request) {
	//just ignore the req for now

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error: while fetching all chirps from db.")
	}

	resJSON := make([]chirpJSON, 0)
	for _, c := range chirps {
		resJSON = append(resJSON, chirpJSON{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	respondWithJSON(w, 200, resJSON)

}
