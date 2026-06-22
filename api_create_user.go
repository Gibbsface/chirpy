package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type createUserRequestJSON struct {
	Email string `json:"email"`
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

	// at this point, we know that reqJSON is valid
	user, err := cfg.db.CreateUser(r.Context(), reqJSON.Email)
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
