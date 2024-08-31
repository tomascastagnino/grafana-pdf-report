package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/services"
)

type PanelHandler struct {
	panelService *services.PanelService
}

func NewPanelHandler(service *services.PanelService) *PanelHandler {
	return &PanelHandler{panelService: service}
}

func (h *PanelHandler) GetPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["dashboard_id"]
	panelID := vars["panel_id"]

	panelURL, err := h.panelService.GetPanel(dashboardID, panelID, *r)
	if err != nil {
		http.Error(w, "Failed to refresh panel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": panelURL})
}
