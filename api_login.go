package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gibbsface/chirpy.git/internal/auth"
	"github.com/Gibbsface/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type loginUserRequestJSON struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginUserResponseJSON struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	//at this point we know the JSON request has been successfully parsed into reqJSON
	//let's extract values into local variables that are easier to read
	email := reqJSON.Email
	password := reqJSON.Password
	expiresInSeconds := 3600 // 1 hr default
	if reqJSON.ExpiresInSeconds != 0 {
		expiresInSeconds = reqJSON.ExpiresInSeconds
	}

	// try to fetch the user obj in db
	user, err := c.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not find user with email %s\n\t%v", email, err))
		return
	}

	// verify their password is correct
	isValid, err := auth.CheckPasswordHash(password, user.Password)
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while checking hash.")
		return
	}

	//create a JWT token for them
	expiresIn := time.Duration(time.Second * time.Duration(expiresInSeconds))
	tokenJWT, err := auth.MakeJWT(user.ID, c.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not make JWT")
	}

	// create access token, store in database
	accessTokenInDB, err := c.db.CreateToken(r.Context(), database.CreateTokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: could not create access token.")
	}

	resJSON := loginUserResponseJSON{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenJWT,
		RefreshToken: accessTokenInDB.Token,
	}

	// fmt.Printf("token generated: %s\n", resJSON.Token)

	// at this point, we know the user was created. Print the results
	// fmt.Printf("User logged in with email %v\n", user.Email)

	//reply with JSON
	respondWithJSON(w, http.StatusOK, resJSON)

}
