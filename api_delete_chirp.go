package main

import (
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
	"github.com/google/uuid"
)

type reqJSON struct {
}

func (c *Config) ApiDeleteChirp(w http.ResponseWriter, r *http.Request) {
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

	// parse chirp id from URL
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error: could not parse chirp ID")
	}

	// fetch chirp id from db
	chirpObj, err := c.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp does not exist")
		return
	}

	if chirpObj.UserID != id {
		respondWithError(w, http.StatusForbidden, "Not your chirp")
		return
	}

	// delete chirp
	err = c.db.DeleteChirp(r.Context(), chirpObj.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while deleting chirp.")
		return
	}

	respondWithError(w, 204, "")
}
