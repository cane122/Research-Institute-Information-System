package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func getUserIDFromRequest(r *http.Request) int {
	userIDStr := r.Header.Get("User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	return userID
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
