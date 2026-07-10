package main

import (
	"net/http"
	"time"

	"github.com/Gibbsface/chirpy.git/internal/auth"
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
	tokenDB, err := c.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not find that token")
		return
	}

	if tokenDB.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "That token was revoked")
		return
	}

	// // create a new access token
	// newToken, err := c.db.CreateToken(r.Context(), database.CreateTokenParams{
	// 	Token:  auth.MakeRefreshToken(),
	// 	UserID: oldToken.UserID,
	// })
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Error: could not create refresh token.")
	// 	return
	// }

	//create a JWT token for them
	expiresIn := time.Duration(3600 * time.Second) // default 1 hour
	jwtToken, err := auth.MakeJWT(tokenDB.UserID, c.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not make JWT")
	}

	resJSON := refreshResponseJSON{
		Token: jwtToken,
	}
	respondWithJSON(w, http.StatusOK, resJSON)

}
