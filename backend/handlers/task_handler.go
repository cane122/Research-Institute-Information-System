package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	taskService      *services.TaskService
	analyticsService *services.AnalyticsService
}

func NewTaskHandler(taskService *services.TaskService, analyticsService *services.AnalyticsService) *TaskHandler {
	return &TaskHandler{
		taskService:      taskService,
		analyticsService: analyticsService,
	}
}

func (h *TaskHandler) GetTasksByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	tasks, err := h.taskService.GetTasksByProject(projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTasksByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	tasks, err := h.taskService.GetTasksByUser(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err := h.taskService.CreateTask(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "TASK_CREATE", "Created task: "+req.NazivZadatka, "task", nil)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskService.GetTaskByID(taskID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.taskService.UpdateTask(taskID, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "TASK_UPDATE", "Updated task", "task", &taskID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.taskService.DeleteTask(taskID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "TASK_DELETE", "Deleted task", "task", &taskID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *TaskHandler) GetTaskComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	comments, err := h.taskService.GetTaskComments(taskID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, comments)
}

func (h *TaskHandler) AddTaskComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID := getUserIDFromRequest(r)
	err = h.taskService.AddTaskComment(taskID, userID, req.Comment)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.analyticsService.LogActivity(&userID, "TASK_COMMENT", "Added comment to task", "task", &taskID)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}
