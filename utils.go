package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt *time.Time
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorRes struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorRes{Error: msg})
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	var lastFetchedAt *time.Time
	if dbFeed.LastFetchedAt.Valid {
		lastFetchedAt = &dbFeed.LastFetchedAt.Time
	}

	return Feed{
		ID:            dbFeed.ID,
		CreatedAt:     dbFeed.CreatedAt,
		UpdatedAt:     dbFeed.UpdatedAt,
		Name:          dbFeed.Name,
		Url:           dbFeed.Url,
		UserID:        dbFeed.UserID,
		LastFetchedAt: lastFetchedAt,
	}
}
