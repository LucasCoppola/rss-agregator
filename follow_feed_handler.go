package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) followFeedsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Feed_id uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode feed params")
		return
	}

	feed, err := apiCfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    params.Feed_id,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}

	respondWithJson(w, http.StatusCreated, feed)
}

func (apiCfg *apiConfig) unfollowFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowId := r.PathValue("feedFollowId")

	if feedFollowId == "" {
		respondWithError(w, http.StatusBadRequest, "Missing id")
		return
	}

	err := apiCfg.DB.UnfollowFeed(r.Context(), database.UnfollowFeedParams{
		ID:     uuid.MustParse(feedFollowId),
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to unfollow feed")
		return
	}

	w.WriteHeader(http.StatusOK)
}
