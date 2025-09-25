package handlers

import (
	"net/http"
	"strconv"

	"github.com/cane/research-institute-system/backend/services"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsService.GetDashboardStats()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetActivityLogs(w http.ResponseWriter, r *http.Request) {
	limit := 50 // Default limit
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.analyticsService.GetActivityLogs(limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}
