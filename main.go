package main

import "net/http"

func main() {
	sMux := http.NewServeMux()
	s := &http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}

	s.ListenAndServe()
}
