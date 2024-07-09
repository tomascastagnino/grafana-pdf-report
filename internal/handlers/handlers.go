package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"fmt"

	"github.com/tomascastagnino/grafana-pdf-report/internal/clients"
)

func HandleReport(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/report/")
    if path == "" {
        http.NotFound(w, r)
        return
    }
    http.ServeFile(w, r, "../../static/index.html")
}

func HandleReportData(w http.ResponseWriter, r *http.Request) {
	dashboardID := strings.TrimPrefix(r.URL.Path, "/api/v1/report/data/")
	dashboardID = strings.TrimSuffix(dashboardID, "/")
	queryParams := r.URL.RawQuery

	if dashboardID == "" {
		http.Error(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	client := clients.GetGrafanaClient()
	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	panelImages := make(map[int]string)
	for _, panel := range dashboard.Panels {
		imageURL := fmt.Sprintf(
			"%s/render/d-solo/%s/?panelId=%d&%s", 
			client.BaseURL,
			dashboardID,
			panel.ID,
			queryParams,
		)
		panelImages[panel.ID] = imageURL
	}

	responseData := struct {
		DashboardID string         `json:"dashboard_id"`
		QueryParams string         `json:"query_params"`
		PanelImages map[int]string `json:"panel_images"`
	}{
		DashboardID: dashboardID,
		QueryParams: queryParams,
		PanelImages: panelImages,
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}