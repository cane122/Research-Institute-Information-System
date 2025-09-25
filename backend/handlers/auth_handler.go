package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService      *services.AuthService
	analyticsService *services.AnalyticsService
}

func NewAuthHandler(authService *services.AuthService, analyticsService *services.AnalyticsService) *AuthHandler {
	return &AuthHandler{
		authService:      authService,
		analyticsService: analyticsService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Login(loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Set session cookie on successful login
	if response.Success && response.User != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    fmt.Sprintf("user_%d", response.User.KorisnikID),
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		// Log activity
		h.analyticsService.LogActivity(&response.User.KorisnikID, "LOGIN", "User logged in", "user", &response.User.KorisnikID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) FirstTimeSetup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username    string `json:"username"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.authService.CompleteFirstTimeSetupByUsername(req.Username, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Log activity
	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "LOGOUT", "User logged out", "user", &userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromRequest(r)
	err := h.authService.ChangePassword(userID, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.analyticsService.LogActivity(&userID, "PASSWORD_CHANGE", "User changed password", "user", &userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	tempPassword, err := h.authService.ResetPassword(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&currentUserID, "PASSWORD_RESET", "Password reset for user", "user", &userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"tempPassword": tempPassword})
}
