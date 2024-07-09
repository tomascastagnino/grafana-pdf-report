package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tomascastagnino/grafana-pdf-report/internal/models"
)

const (
	GrafanaURL = "http://localhost:3000" 
	ApiKey = "_"
	dashboardURL = "/api/dashboards/uid/"
)

type grafanaClient struct {
	http.Client

	BaseURL string
	dashboardURL string
	ApiKey string
}

func GetGrafanaClient() *grafanaClient {
	return &grafanaClient{
		BaseURL: GrafanaURL,
		dashboardURL: dashboardURL,
		ApiKey: ApiKey,
	}
}

func (c *grafanaClient) GetDashboard(dashboardID string) (*models.Dashboard, error) {
	url := c.BaseURL+c.dashboardURL+dashboardID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get dashboard: %s", string(body))
	}

	var result struct {
		Dashboard models.Dashboard `json:"dashboard"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Dashboard, nil
}

// func (c *grafanaClient) BuildRequest(method string, url string, key string) {

// }

// func (c *grafanaClient) GetPanelsURL() () {

// }

// func (c *grafanaClient) GetPanels(panelsURL []string) () {

// }