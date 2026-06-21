package main

import (
	"fmt"
	"net/http"
)

func (cfg *Config) AdminReset(w http.ResponseWriter, r *http.Request) {

	//verify
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Error: forbidden")
		return
	}

	cfg.fileserverHits.Swap(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		fmt.Printf("Error while dropping users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error while dropping users")
		return
	}

	w.WriteHeader(200)
}
