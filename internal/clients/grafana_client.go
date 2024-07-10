package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"image"
	"image/png"
    "io"
    "net/http"
	"path/filepath"
    "os"
	"log"
	"strconv"
	"sync"

	"github.com/nfnt/resize"
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
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for _, panel := range dashboard.Panels {
		if panel.Tag == "remove" {
			continue	
		}
		if panel.Type == "text" {
			mu.Lock()
			panels[panel.ID] = models.Panel{
				ID:      panel.ID,
				Type:    panel.Type,
				GridPos: panel.GridPos,
				Options: panel.Options,
				Tag: panel.Tag,
			}
			mu.Unlock()
			continue
		}
		wg.Add(1)
		go func(panel models.Panel) {
			defer wg.Done()

			url := getImageURL(c.BaseURL, uid, panel.ID, params) 
			localImagePath := filepath.Join("../../static/images", filepath.Base(strconv.Itoa(panel.ID))) + ".png"
			err := c.downloadImage(url, localImagePath, panel.GridPos)
			if err != nil {
				log.Printf("Failed to download image for panel %d: %v", panel.ID, err)
				return
			}

			mu.Lock()
			panels[panel.ID] = models.Panel{
				ID: panel.ID,
				URL: "/static/images/" + strconv.Itoa(panel.ID) + ".png",
				GridPos: panel.GridPos,
				Tag: panel.Tag,
			}
			mu.Unlock()
		}(panel)
	}
	wg.Wait()
	return panels
}

func getImageURL (b string, dashboard string, panel int, params string) string {
	return fmt.Sprintf("%s/render/d-solo/%s/?panelId=%d&%s", b, dashboard, panel, params)
}

func (c *grafanaClient) downloadImage(url, filePath string, pos models.GridPos) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()	

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return err
	}

	// Width and height for 24 columns
	baseWidth := 2400.0
	colWidth := baseWidth / 24.0
	rowHeight := colWidth * 0.5625  // 16:9 aspect ratio
	targetWidth := uint(colWidth * float64(pos.W))
	targetHeight := uint(rowHeight * float64(pos.H))
	
	resizedImg := resize.Resize(targetWidth, targetHeight, img, resize.Lanczos3)
	file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
	png.Encode(file, resizedImg)

    _, err = io.Copy(file, resp.Body)
    return err
}