package handlers

import "net/http"

type statusRes struct {
	Status string `json:"status"`
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, statusRes{Status: "ok"})
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}
