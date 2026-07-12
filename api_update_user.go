package main

import (
	"encoding/json"
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
	"github.com/Gibbsface/chirpy.git/internal/database"
)

type requestJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseJSON struct {
	Email string `json:"email"`
}

func (c *Config) ApiUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqJSON := &requestJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	//validate JWT
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not parse Bearer Token")
		return
	}
	id, err := auth.ValidateJWT(jwt, c.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT is invalid")
		return
	}

	// hash password before storing in db
	hashedPassword, err := auth.HashPassword(reqJSON.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	//query
	userObj, err := c.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:    reqJSON.Email,
		Password: hashedPassword,
		ID:       id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not update user in db.")
	}

	respondWithJSON(w, http.StatusOK, responseJSON{
		Email: userObj.Email,
	})
}
