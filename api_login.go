package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
)

type loginUserRequestJSON struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (c *Config) ApiLogin(w http.ResponseWriter, r *http.Request) {
	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := &loginUserRequestJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	user, err := c.db.GetUserByEmail(r.Context(), reqJSON.Email)

	isValid, err := auth.CheckPasswordHash(reqJSON.Password, user.Password)
	if isValid {
		resJSON := userJSON{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		// at this point, we know the user was created. Print the results
		fmt.Printf("User logged in with email %v\n", user.Email)

		//reply with JSON
		respondWithJSON(w, http.StatusOK, resJSON)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
	}
}
