package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"time"
	"log"


	"github.com/tomascastagnino/grafana-pdf-report/internal/clients"
	"github.com/tomascastagnino/grafana-pdf-report/internal/models"
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
	startTime := time.Now()
	dashboardID := strings.TrimPrefix(r.URL.Path, "/api/v1/report/data/")
	dashboardID = strings.TrimSuffix(dashboardID, "/")
	queryParams := r.URL.RawQuery

	if dashboardID == "" {
		http.Error(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	client := clients.GetGrafanaClient()
	dashboardStartTime := time.Now()
	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	DeleteImages("../../static/images")
	log.Printf("Time taken to get dashboard: %v", time.Since(dashboardStartTime))

	// Measure time to get panels
	panelsStartTime := time.Now()
	panels := client.GetPanels(*dashboard, dashboardID, queryParams)
	log.Printf("Time taken to get panels: %v", time.Since(panelsStartTime))
	responseData := struct {
		DashboardID string              `json:"dashboard_id"`
		QueryParams string              `json:"query_params"`
		Panels      map[int]models.Panel `json:"panels"`
	}{
		DashboardID: dashboardID,
		QueryParams: queryParams,
		Panels: panels,
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	log.Printf("Total time taken: %v", time.Since(startTime))
}