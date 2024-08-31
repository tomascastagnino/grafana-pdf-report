package clients

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/utils"
)

type grafanaAPIClient struct {
	baseURL    string
	httpClient *http.Client
	headers    http.Header
}

func NewGrafanaAPIClient(baseURL string, httpClient *http.Client, headers http.Header) GrafanaClient {
	return &grafanaAPIClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		headers:    headers,
	}
}

// GetDashboard fetches the dashboard data from Grafana.
func (g *grafanaAPIClient) GetDashboard(dashboardID string, h http.Header) (*models.Dashboard, error) {
	url := g.baseURL + dashboardID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = h

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get dashboard: %s", string(body))
	}

	var result struct {
		Dashboard models.Dashboard `json:"dashboard"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Dashboard, nil
}

func (g *grafanaAPIClient) GetPanelImage(dID string, params url.Values, h http.Header) (string, error) {
	url := fmt.Sprintf("%s/%s?%s", internal.ImageRendererURL, dID, params.Encode())
	in := utils.ImgName(params)
	path := filepath.Join(internal.ImageDir, in)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for image: %w", err)
	}
	req.Header = h

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("failed to encode image to PNG: %w", err)
	}

	return filepath.Join(internal.WebImageDir, filepath.Base(in)), nil
}

// GetAllDashboards fetches all dashboards from Grafana.
func (g *grafanaAPIClient) GetAllDashboards() ([]models.Dashboard, error) {
	panic("TODO: To be implemented in the future.")
}
