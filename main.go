package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Gibbsface/chirpy.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening db connection: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	// initialize server config state
	cfg := Config{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
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
