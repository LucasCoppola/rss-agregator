package main

import "net/http"

type statusRes struct {
	Status string `json:"status"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, statusRes{Status: "ok"})
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
