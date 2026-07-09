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

type createUserRequestJSON struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type createUserResponseJSON struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	// Token     string    `json:"token"`
}

func (c *Config) ApiCreateUser(w http.ResponseWriter, r *http.Request) {
	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := &createUserRequestJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	//at this point we know the JSON request has been successfully parsed into reqJSON
	//let's extract values into local variables that are easier to read
	email := reqJSON.Email
	password := reqJSON.Password
	// expiresInSeconds := 3600 // default

	//hash the pw
	hash, err := auth.HashPassword(password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing the password.")
	}

	// at this point, we know that reqJSON is valid
	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    email,
		Password: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	//create a JWT token for them
	// token, err := auth.MakeJWT(user.ID, c.secret, time.Duration(time.Second*time.Duration(expiresInSeconds)))
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Error: could not make JWT")
	// }

	resJSON := createUserResponseJSON{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		// Token:     token,
	}

	// at this point, we know the user was created. Print the results
	fmt.Printf("User created with email %v\n", user.Email)

	//reply with JSON
	respondWithJSON(w, http.StatusCreated, resJSON)
}
