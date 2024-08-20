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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	dashboardID := params.Get("dashboardId")

	if dashboardID == "" {
		http.Error(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	client := clients.GetGrafanaClient(&r.Header)
	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client.DeleteImages("../../static/images")
	panels := client.GetPanels(*dashboard, *r)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	c := clients.GetGrafanaClient(&r.Header)
	imageURL, err := c.GetRefreshedPanelURL(*r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse := map[string]string{"url": imageURL}
	json.NewEncoder(w).Encode(jsonResponse)
}
