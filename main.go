package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/GevaYo/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	JWTSecret      string
}

func main() {
	godotenv.Load()
	const filePathRoot = "."
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	JWTSecret := os.Getenv("JWT_SECRET")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to connect to the DB: %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		JWTSecret:      JWTSecret,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", (http.FileServer(http.Dir(".")))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.validateChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.chirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.chirpHandler)

	mux.HandleFunc("POST /api/users", apiCfg.usersHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
