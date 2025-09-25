// ============================================================================
// task_comment_handler.go - Task Comments HTTP Handlers
// ============================================================================

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type TaskCommentHandler struct {
	taskCommentService *services.TaskCommentService
	analyticsService   *services.AnalyticsService
}

func NewTaskCommentHandler(taskCommentService *services.TaskCommentService, analyticsService *services.AnalyticsService) *TaskCommentHandler {
	return &TaskCommentHandler{
		taskCommentService: taskCommentService,
		analyticsService:   analyticsService,
	}
}

// GetTaskComments retrieves all comments for a specific task
// GET /tasks/{id}/comments
func (h *TaskCommentHandler) GetTaskComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	comments, err := h.taskCommentService.GetCommentsByTask(taskID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log activity
	userID := getUserIDFromRequest(r)
	if userID > 0 {
		h.analyticsService.LogActivity(&userID, "READ", "Retrieved task comments", "task_comments", &taskID)
	}

	respondWithJSON(w, http.StatusOK, comments)
}

// CreateTaskComment adds a new comment to a task
// POST /tasks/{id}/comments
func (h *TaskCommentHandler) CreateTaskComment(w http.ResponseWriter, r *http.Request) {
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
		TekstKomentara string `json:"tekst_komentara"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if requestData.TekstKomentara == "" {
		respondWithError(w, http.StatusBadRequest, "Comment text is required")
		return
	}

	comment, err := h.taskCommentService.CreateComment(taskID, userID, requestData.TekstKomentara)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log activity
	h.analyticsService.LogActivity(&userID, "CREATE", "Created task comment", "task_comment", &comment.KomentarID)

	respondWithJSON(w, http.StatusCreated, comment)
}

// UpdateTaskComment updates an existing comment
// PUT /comments/{id}
func (h *TaskCommentHandler) UpdateTaskComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var requestData struct {
		TekstKomentara string `json:"tekst_komentara"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if requestData.TekstKomentara == "" {
		respondWithError(w, http.StatusBadRequest, "Comment text is required")
		return
	}

	// Check if user is admin (you might want to implement a proper role check)
	isAdmin := h.isUserAdmin(r)

	err = h.taskCommentService.UpdateComment(commentID, userID, requestData.TekstKomentara, isAdmin)
	if err != nil {
		if err.Error() == "unauthorized: can only edit own comments" {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else if err.Error() == "comment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Log activity
	h.analyticsService.LogActivity(&userID, "UPDATE", "Updated task comment", "task_comment", &commentID)

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Comment updated successfully"})
}

// DeleteTaskComment removes a comment
// DELETE /comments/{id}
func (h *TaskCommentHandler) DeleteTaskComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user is admin
	isAdmin := h.isUserAdmin(r)

	err = h.taskCommentService.DeleteComment(commentID, userID, isAdmin)
	if err != nil {
		if err.Error() == "unauthorized: can only delete own comments" {
			respondWithError(w, http.StatusForbidden, err.Error())
		} else if err.Error() == "comment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Log activity
	h.analyticsService.LogActivity(&userID, "DELETE", "Deleted task comment", "task_comment", &commentID)

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully"})
}

// GetUserComments retrieves all comments made by a specific user
// GET /users/{id}/comments
func (h *TaskCommentHandler) GetUserComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Only allow users to view their own comments or admins to view any
	requestingUserID := getUserIDFromRequest(r)
	if requestingUserID != userID && !h.isUserAdmin(r) {
		respondWithError(w, http.StatusForbidden, "Unauthorized to view these comments")
		return
	}

	comments, err := h.taskCommentService.GetCommentsByUser(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log activity
	if requestingUserID > 0 {
		h.analyticsService.LogActivity(&requestingUserID, "READ", "Retrieved user comments", "user_comments", &userID)
	}

	respondWithJSON(w, http.StatusOK, comments)
}

// GetCommentByID retrieves a specific comment by ID
// GET /comments/{id}
func (h *TaskCommentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := h.taskCommentService.GetCommentByID(commentID)
	if err != nil {
		if err.Error() == "comment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Log activity
	userID := getUserIDFromRequest(r)
	if userID > 0 {
		h.analyticsService.LogActivity(&userID, "READ", "Retrieved comment", "task_comment", &commentID)
	}

	respondWithJSON(w, http.StatusOK, comment)
}

// Helper function to check if user is admin
// This is a simplified version - you should implement proper role checking
func (h *TaskCommentHandler) isUserAdmin(r *http.Request) bool {
	// Get role from request header or context
	role := r.Header.Get("User-Role")
	return role == "Administrator" || role == "Admin"
}