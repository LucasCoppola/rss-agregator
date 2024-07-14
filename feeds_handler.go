package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createFeedsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type feedParams struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := feedParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode feed params")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	followFeed, err := apiCfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}

	type feedResponse struct {
		Feed       Feed              `json:"feed"`
		FeedFollow database.FeedUser `json:"feed_follow"`
	}

	parsedFeed := databaseFeedToFeed(feed)

	respondWithJson(w, http.StatusCreated, feedResponse{Feed: parsedFeed, FeedFollow: followFeed})
}

func (apiCfg *apiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get feeds")
		return
	}

	var parsedFeeds []Feed

	for _, feed := range feeds {
		parsedFeeds = append(parsedFeeds, databaseFeedToFeed(feed))
	}

	respondWithJson(w, http.StatusOK, parsedFeeds)
}
