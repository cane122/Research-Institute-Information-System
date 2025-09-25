package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/services"
	"github.com/gorilla/mux"
)

type WorkflowHandler struct {
	workflowService  *services.WorkflowService
	analyticsService *services.AnalyticsService
}

func NewWorkflowHandler(workflowService *services.WorkflowService, analyticsService *services.AnalyticsService) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService:  workflowService,
		analyticsService: analyticsService,
	}
}

func (h *WorkflowHandler) GetAllWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := h.workflowService.GetAllWorkflows()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workflows)
}

func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var workflow models.RadniTokovi
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err := h.workflowService.CreateWorkflow(workflow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "WORKFLOW_CREATE", "Created workflow: "+workflow.Naziv, "workflow", nil)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

func (h *WorkflowHandler) GetWorkflowPhases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workflowID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	phases, err := h.workflowService.GetWorkflowPhases(workflowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, phases)
}

func (h *WorkflowHandler) CreatePhase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workflowID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	var phase models.Faze
	if err := json.NewDecoder(r.Body).Decode(&phase); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	phase.RadniTokID = workflowID
	err = h.workflowService.CreatePhase(phase)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := getUserIDFromRequest(r)
	h.analyticsService.LogActivity(&userID, "WORKFLOW_PHASE_CREATE", "Created phase: "+phase.NazivFaze, "workflow", &workflowID)

	respondWithJSON(w, http.StatusCreated, map[string]bool{"success": true})
}
