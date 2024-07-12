package main

import (
	"net/http"
	"strings"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		apiKey := strings.TrimPrefix(header, "ApiKey ")

		if apiKey == "" {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		handler(w, r, user)
	}
}
