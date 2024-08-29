package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
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
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Bad Request: Invalid query parameters", http.StatusBadRequest)
		return
	}

	dashboardID := params.Get("dashboardId")
	if dashboardID == "" {
		http.Error(w, "Bad Request: Missing dashboard ID", http.StatusBadRequest)
		return
	}

	client := clients.GetGrafanaClient(&r.Header)

	// Fetch the dashboard information
	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		http.Error(w, "Failed to get dashboard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Clean up images
	err = client.DeleteImages("../../static/images")
	if err != nil {
		http.Error(w, "Failed to delete previous images: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch and process panels
	panels, err := client.GetPanels(*dashboard, *r)
	if err != nil {
		http.Error(w, "Failed to fetch panel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := struct {
		DashboardID string               `json:"dashboard_id"`
		QueryParams url.Values           `json:"query_params"`
		Panels      map[int]models.Panel `json:"panels"`
	}{
		DashboardID: dashboardID,
		QueryParams: params,
		Panels:      panels,
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal response data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	c := clients.GetGrafanaClient(&r.Header)
	imageURL, err := c.GetRefreshedPanelURL(*r)
	if err != nil {
		http.Error(w, "Failed to refresh image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse := map[string]string{"url": imageURL}
	json.NewEncoder(w).Encode(jsonResponse)
}
