package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
    "io"
    "net/http"
	"path/filepath"
    "os"
	"log"
	"strconv"

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

func (c *grafanaClient) GetPanels(dashboard models.Dashboard, uid string, params string) map[int]models.Panel {
	panels := make(map[int]models.Panel)
	for _, panel := range dashboard.Panels {
		url := getImageURL(c.BaseURL, uid, panel.ID, params) 
		localImagePath := filepath.Join("../../static/images", filepath.Base(strconv.Itoa(panel.ID)))
		localImagePath = localImagePath+".png"
		c.downloadImage(url, localImagePath)
		panels[panel.ID] = models.Panel{
			ID: panel.ID,
			URL: "/static/images/"+strconv.Itoa(panel.ID)+".png",
			GridPos: panel.GridPos,
		}
	}
	return panels
}

func getImageURL (b string, dashboard string, panel int, params string) string {
	url := fmt.Sprintf("%s/render/d-solo/%s/?panelId=%d&%s", b, dashboard, panel, params)
	return url
}

func (c *grafanaClient) downloadImage(url, filePath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.Do(req)
	if err != nil {
		log.Printf("hello", err)
		return err
	}
	defer resp.Body.Close()	
	file, err := os.Create(filePath)
    if err != nil {
		log.Printf("hello", err)
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    return err
}