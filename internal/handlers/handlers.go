package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

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
	dashboardID := strings.TrimPrefix(r.URL.Path, "/api/v1/report/data/")
	dashboardID = strings.TrimSuffix(dashboardID, "/")

	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

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
	client.DeleteImages("../../static/images")
	panels := client.GetPanels(*dashboard, dashboardID, queryParams)
	responseData := struct {
		DashboardID string               `json:"dashboard_id"`
		QueryParams url.Values           `json:"query_params"`
		Panels      map[int]models.Panel `json:"panels"`
	}{
		DashboardID: dashboardID,
		QueryParams: queryParams,
		Panels:      panels,
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
