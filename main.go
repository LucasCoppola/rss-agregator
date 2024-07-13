package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	DB_URL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.getUserHandler))
	mux.HandleFunc("POST /v1/users", apiCfg.createUserHandler)

	mux.HandleFunc("GET /v1/feeds", apiCfg.getFeedsHandler)
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.createFeedsHandler))

	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.getFollowedFeedsHandler))
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.followFeedHandler))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.unfollowFeedHandler))

	mux.HandleFunc("GET /v1/healthz", healthzHandler)
	mux.HandleFunc("GET /v1/err", errorHandler)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
