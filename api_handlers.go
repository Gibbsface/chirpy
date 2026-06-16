package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func ApiHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func ApiValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "")
		return
	}

	type returnVals struct {
		CreatedAt time.Time `json:"created_at"`
		ID        int       `json:"id"`
	}
	respBody := returnVals{
		CreatedAt: time.Now(),
		ID:        123,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
