package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectService   *services.ProjectService
	analyticsService *services.AnalyticsService
}

func NewProjectHandler(projectService *services.ProjectService, analyticsService *services.AnalyticsService) *ProjectHandler {
	return &ProjectHandler{
		projectService:   projectService,
		analyticsService: analyticsService,
	}
}

func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetAllProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err := h.projectService.CreateProject(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "PROJECT_CREATE", "Created project: "+req.NazivProjekta, "project", nil)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	respondWithJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var project models.Projekti
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.projectService.UpdateProject(projectID, project)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "PROJECT_UPDATE", "Updated project", "project", &projectID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	err = h.projectService.DeleteProject(projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "PROJECT_DELETE", "Deleted project", "project", &projectID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *ProjectHandler) GetProjectMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	members, err := h.projectService.GetProjectMembers(projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, members)
}

func (h *ProjectHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req struct {
		UserID int `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.projectService.AddProjectMember(projectID, req.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "PROJECT_MEMBER_ADD", "Added member to project", "project", &projectID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *ProjectHandler) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	memberUserID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.projectService.RemoveProjectMember(projectID, memberUserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "PROJECT_MEMBER_REMOVE", "Removed member from project", "project", &projectID)

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}
