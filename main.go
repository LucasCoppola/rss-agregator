package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LucasCoppola/rss-aggregator/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", handlers.HealthzHandler)
	mux.HandleFunc("GET /v1/err", handlers.ErrorHandler)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}
