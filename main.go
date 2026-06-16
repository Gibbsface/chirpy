package main

import (
	"net/http"
	"sync/atomic"
)

// func (cfg *Config) middlewareMetricsInc(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		next.ServeHTTP(w, r)
// 	})
// }

func (cfg *Config) AppHandler() http.Handler {
	cfg.fileserverHits.Add(1)
	return http.StripPrefix("/app", http.FileServer(http.Dir(".")))
}

func main() {
	sMux := http.NewServeMux()
	cfg := Config{
		fileserverHits: atomic.Int32{},
	}

	sMux.HandleFunc("GET /api/healthz", ApiHealthz)
	sMux.HandleFunc("POST /api/validate_chirp", ApiValidateChirp)

	sMux.HandleFunc("GET /admin/metrics", cfg.AdminMetrics)
	sMux.HandleFunc("POST /admin/reset", cfg.AdminReset)

	sMux.Handle("/app/", cfg.AppHandler())

	s := &http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}

	s.ListenAndServe()

}
