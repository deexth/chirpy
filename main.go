package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/deexth/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("The db URL is missing")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("No platform set")
	}

	tSecret := os.Getenv("TOKEN_SECRET")
	if tSecret == "" {
		log.Fatal("No token secret found")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect to the db: %v", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	apicfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		tokenSecret:    tSecret,
	}
	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Fatalf("couldn't write to the body: %v", err)
		}
	})
	mux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handleReset)
	mux.HandleFunc("POST /api/users", apicfg.handleUsers)
	mux.HandleFunc("POST /api/login", apicfg.handleLogin)
	mux.HandleFunc("POST /api/chirps", apicfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apicfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.handleGetChirp)
	mux.HandleFunc("POST /api/refresh", apicfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apicfg.handleRevoke)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("the server couldn't start: %v\n", err)
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
