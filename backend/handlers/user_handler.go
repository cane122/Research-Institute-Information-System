package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService      *services.UserService
	analyticsService *services.AnalyticsService
}

func NewUserHandler(userService *services.UserService, analyticsService *services.AnalyticsService) *UserHandler {
	return &UserHandler{
		userService:      userService,
		analyticsService: analyticsService,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User         models.Korisnici `json:"user"`
		TempPassword string           `json:"tempPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.userService.CreateUser(req.User, req.TempPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&currentUserID, "USER_CREATE", "Created new user: "+req.User.KorisnickoIme, "user", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.Korisnici
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateUser(userID, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&currentUserID, "USER_UPDATE", "Updated user", "user", &userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&currentUserID, "USER_DELETE", "Deleted user", "user", &userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *UserHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.userService.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}
