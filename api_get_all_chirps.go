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

	resJSON := make([]createChirpResponseJSON, 0)
	for _, c := range chirps {
		resJSON = append(resJSON, createChirpResponseJSON{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	respondWithJSON(w, 200, resJSON)

}
