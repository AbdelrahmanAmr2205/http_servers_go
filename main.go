package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldn't load environment variables:", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET_KEY")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	if secretKey == "" {
		log.Fatal("SECRET_KEY must be set")
	}

	db, err := sql.Open("postgres", dbURL)

	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{db: database.New(db), platform: platform, secretKey: secretKey}

	sMux := http.NewServeMux()
	f := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	sMux.Handle("GET /app/", cfg.middlewareMetricsInc(f))
	sMux.HandleFunc("GET /api/healthz", handleHealthz)
	sMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	sMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	sMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	sMux.Handle("POST /api/chirps", cfg.middlewareAuth(cfg.handlerCreateChirp))
	sMux.HandleFunc("GET /api/chirps", cfg.getAllChirps)
	sMux.HandleFunc("GET /api/chirps/{id}", cfg.getChirp)
	sMux.HandleFunc("POST /api/login", cfg.handlerLogin)
	sMux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	sMux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	sMux.Handle("PUT /api/users", cfg.middlewareAuth(cfg.handlerEditUser))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
