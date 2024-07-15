package main

import (
	"net/http"
	"strconv"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
)

const DEFAULT_LIMIT = 10

func (apiCfg *apiConfig) getPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	var limit int32
	limitQuery := r.URL.Query().Get("limit")

	if limitQuery != "" {
		limitInt, err := strconv.Atoi(limitQuery)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse limit")
			return
		}

		limit = int32(limitInt)
	} else {
		limit = DEFAULT_LIMIT
	}

	posts, err := apiCfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  limit,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}

	respondWithJson(w, http.StatusOK, posts)
}
