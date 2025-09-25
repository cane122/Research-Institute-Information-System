package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

var db *sql.DB

func setDB(database *sql.DB) {
	db = database
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if sessionCookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse user ID from session cookie (format: "user_{id}")
		sessionValue := sessionCookie.Value
		if !strings.HasPrefix(sessionValue, "user_") {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		userIDStr := strings.TrimPrefix(sessionValue, "user_")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Get user role from database
		var role string
		query := `
			SELECT u.naziv_uloge 
			FROM Korisnici k 
			JOIN Uloge u ON k.uloga_id = u.uloga_id 
			WHERE k.korisnik_id = $1 AND k.status = 'aktivan'
		`
		err = db.QueryRow(query, userID).Scan(&role)
		if err != nil {
			http.Error(w, "User not found or inactive", http.StatusUnauthorized)
			return
		}

		// Set headers for downstream handlers
		r.Header.Set("User-ID", strconv.Itoa(userID))
		r.Header.Set("User-Role", role)

		next.ServeHTTP(w, r)
	})
}
