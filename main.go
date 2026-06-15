package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	sMux := http.NewServeMux()
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	sMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	sMux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load())
	})

	sMux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)
		w.WriteHeader(200)
	})

	sMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	s := &http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}

	s.ListenAndServe()

}
