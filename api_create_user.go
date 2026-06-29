package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gibbsface/chirpy.git/internal/auth"
	"github.com/Gibbsface/chirpy.git/internal/database"
)

type createUserRequestJSON struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *Config) ApiCreateUser(w http.ResponseWriter, r *http.Request) {
	// attempt to decode JSON from request
	decoder := json.NewDecoder(r.Body)
	reqJSON := &createUserRequestJSON{}
	err := decoder.Decode(reqJSON)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read JSON request")
		return
	}

	//hash the pw
	hash, err := auth.HashPassword(reqJSON.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing the password.")
	}

	// at this point, we know that reqJSON is valid
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    reqJSON.Email,
		Password: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	resJSON := userJSON{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	// at this point, we know the user was created. Print the results
	fmt.Printf("User created with email %v\n", user.Email)

	//reply with JSON
	respondWithJSON(w, http.StatusCreated, resJSON)
}
