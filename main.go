package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	// initialize server config state
	cfg := Config{
		fileserverHits: atomic.Int32{},
	}

	//Create serverMux and register handlers
	sMux := http.NewServeMux()

	sMux.HandleFunc("GET /api/healthz", ApiHealthz)
	sMux.HandleFunc("POST /api/validate_chirp", ApiValidateChirp)

	sMux.HandleFunc("GET /admin/metrics", cfg.AdminMetrics)
	sMux.HandleFunc("POST /admin/reset", cfg.AdminReset)

	// file server handler
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	sMux.Handle("/app/", cfg.MiddlewareMetricsInc(fsHandler))

	// initialize server struct and start serving
	s := &http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}

	s.ListenAndServe()

}
