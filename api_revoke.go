package main

import (
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
)

func (c *Config) ApiRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not detect Bearer token")
		return
	}

	//get refresh token in db
	_, err = c.db.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not find that token")
		return
	}

	respondWithError(w, http.StatusNoContent, "")
}
