// ============================================================================
// phase_change_request_handler.go - Phase Change Requests HTTP Handlers
// ============================================================================

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type PhaseChangeRequestHandler struct {
	phaseChangeService *services.PhaseChangeRequestService
	analyticsService   *services.AnalyticsService
}

func NewPhaseChangeRequestHandler(phaseChangeService *services.PhaseChangeRequestService, analyticsService *services.AnalyticsService) *PhaseChangeRequestHandler {
	return &PhaseChangeRequestHandler{
		phaseChangeService: phaseChangeService,
		analyticsService:   analyticsService,
	}
}

// GetTaskPhaseChangeRequests retrieves all phase change requests for a specific task
// GET /tasks/{id}/phase-requests
func (h *PhaseChangeRequestHandler) GetTaskPhaseChangeRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	requests, err := h.phaseChangeService.GetRequestsByTask(taskID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log activity
	userID := getUserIDFromRequest(r)
	if userID > 0 {
		h.analyticsService.LogActivity(&userID, "READ", "Retrieved task phase change requests", "phase_change_requests", &taskID)
	}

	respondWithJSON(w, http.StatusOK, requests)
}

// CreatePhaseChangeRequest creates a new phase change request
// POST /tasks/{id}/phase-requests
func (h *PhaseChangeRequestHandler) CreatePhaseChangeRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var requestData struct {
		ZahtevanaFazaID int    `json:"zahtevana_faza_id"`
		Komentar        string `json:"komentar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if requestData.ZahtevanaFazaID == 0 {
		respondWithError(w, http.StatusBadRequest, "Requested phase ID is required")
		return
	}

	request, err := h.phaseChangeService.CreateRequest(taskID, userID, requestData.ZahtevanaFazaID, requestData.Komentar)
	if err != nil {
		if err.Error() == "there is already a pending phase change request for this task" {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Log activity
	h.analyticsService.LogActivity(&userID, "CREATE", "Created phase change request", "phase_change_request", &request.ZahtevID)

	respondWithJSON(w, http.StatusCreated, request)
}

// UpdatePhaseChangeRequestStatus approves or rejects a phase change request
// PUT /phase-requests/{id}/status
func (h *PhaseChangeRequestHandler) UpdatePhaseChangeRequestStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user has permission to approve/reject requests (admin or project manager)
	if !h.canManagePhaseRequests(r) {
		respondWithError(w, http.StatusForbidden, "Insufficient permissions to manage phase change requests")
		return
	}

	var requestData struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if requestData.Status != "Odobren" && requestData.Status != "Odbijen" {
		respondWithError(w, http.StatusBadRequest, "Status must be 'Odobren' or 'Odbijen'")
		return
	}

	err = h.phaseChangeService.UpdateRequestStatus(requestID, requestData.Status, userID)
	if err != nil {
		if err.Error() == "phase change request not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "request has already been processed" {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Log activity
	h.analyticsService.LogActivity(&userID, "UPDATE", "Updated phase change request status: "+requestData.Status, "phase_change_request", &requestID)

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Request status updated successfully"})
}

// GetPendingPhaseChangeRequests retrieves all pending phase change requests
// GET /phase-requests/pending
func (h *PhaseChangeRequestHandler) GetPendingPhaseChangeRequests(w http.ResponseWriter, r *http.Request) {
	// Only allow users with management permissions to view all pending requests
	if !h.canManagePhaseRequests(r) {
		respondWithError(w, http.StatusForbidden, "Insufficient permissions to view pending requests")
		return
	}

	requests, err := h.phaseChangeService.GetPendingRequests()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log activity
	userID := getUserIDFromRequest(r)
	if userID > 0 {
		h.analyticsService.LogActivity(&userID, "READ", "Retrieved pending phase change requests", "pending_phase_requests", nil)
	}

	respondWithJSON(w, http.StatusOK, requests)
}

// GetUserPhaseChangeRequests retrieves all phase change requests made by a specific user
// GET /users/{id}/phase-requests
func (h *PhaseChangeRequestHandler) GetUserPhaseChangeRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Only allow users to view their own requests or admins/managers to view any
	requestingUserID := getUserIDFromRequest(r)
	if requestingUserID != userID && !h.canManagePhaseRequests(r) {
		respondWithError(w, http.StatusForbidden, "Unauthorized to view these requests")
		return
	}

	requests, err := h.phaseChangeService.GetRequestsByUser(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log activity
	if requestingUserID > 0 {
		h.analyticsService.LogActivity(&requestingUserID, "READ", "Retrieved user phase change requests", "user_phase_requests", &userID)
	}

	respondWithJSON(w, http.StatusOK, requests)
}

// GetPhaseChangeRequestByID retrieves a specific phase change request by ID
// GET /phase-requests/{id}
func (h *PhaseChangeRequestHandler) GetPhaseChangeRequestByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	request, err := h.phaseChangeService.GetRequestByID(requestID)
	if err != nil {
		if err.Error() == "phase change request not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Check if user has permission to view this request
	userID := getUserIDFromRequest(r)
	if userID != request.PodnosilacZahtevaID && !h.canManagePhaseRequests(r) {
		respondWithError(w, http.StatusForbidden, "Unauthorized to view this request")
		return
	}

	// Log activity
	if userID > 0 {
		h.analyticsService.LogActivity(&userID, "READ", "Retrieved phase change request", "phase_change_request", &requestID)
	}

	respondWithJSON(w, http.StatusOK, request)
}

// DeletePhaseChangeRequest removes a phase change request
// DELETE /phase-requests/{id}
func (h *PhaseChangeRequestHandler) DeletePhaseChangeRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user is admin or project manager
	isAdmin := h.canManagePhaseRequests(r)

	err = h.phaseChangeService.DeleteRequest(requestID, userID, isAdmin)
	if err != nil {
		if err.Error() == "unauthorized: can only delete own requests" {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else if err.Error() == "phase change request not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "can only delete pending requests" {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Log activity
	h.analyticsService.LogActivity(&userID, "DELETE", "Deleted phase change request", "phase_change_request", &requestID)

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Request deleted successfully"})
}

// Helper function to check if user can manage phase requests (admin or project manager)
func (h *PhaseChangeRequestHandler) canManagePhaseRequests(r *http.Request) bool {
	role := r.Header.Get("User-Role")
	return role == "Administrator" || role == "Admin" || role == "Rukovodilac projekta" || role == "Project Manager"
}