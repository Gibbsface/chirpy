package main

import (
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
	"github.com/Gibbsface/chirpy.git/internal/database"
)

type refreshResponseJSON struct {
	Token string `json:"token"`
}

func (c *Config) ApiRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not detect Bearer token")
		return
	}

	//get refresh token in db
	oldToken, err := c.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not find that token")
		return
	}

	// create a new access token
	newToken, err := c.db.CreateToken(r.Context(), database.CreateTokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: oldToken.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not create refresh token.")
		return
	}

	resJSON := refreshResponseJSON{
		Token: newToken.Token,
	}
	respondWithJSON(w, http.StatusOK, resJSON)

}
