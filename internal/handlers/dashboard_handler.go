package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/services"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: service}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["dashboard_id"]

	dashboard, err := h.dashboardService.GetDashboard(dashboardID, r)
	if err != nil {
		http.Error(w, "Failed to get dashboard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (h *DashboardHandler) ListDashboards(w http.ResponseWriter, r *http.Request) {
	dashboards, err := h.dashboardService.ListDashboards(r)
	if err != nil {
		http.Error(w, "Failed to list dashboards: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboards)
}
