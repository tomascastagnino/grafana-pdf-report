package clients

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/tomascastagnino/grafana-pdf-report/internal/models"
	"gopkg.in/ini.v1"
)

var (
	GrafanaURL   string
	dashboardURL = "/api/dashboards/uid/"
)

type grafanaClient struct {
	http.Client

	BaseURL      string
	dashboardURL string
	Header       http.Header
}

func GetGrafanaClient(r *http.Header) *grafanaClient {
	return &grafanaClient{
		BaseURL:      getGrafanaURL(),
		dashboardURL: dashboardURL,
		Header:       r.Clone(),
	}
}

func getGrafanaURL() string {
	cfg, err := ini.Load("../../config.ini")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	url := cfg.Section("server").Key("GrafanaURL").String()
	if url == "" {
		log.Fatalf("GrafanaURL not set in config file")
	}
	return url
}

func (c *grafanaClient) GetDashboard(dashboardID string) (*models.Dashboard, error) {
	url := c.BaseURL + c.dashboardURL + dashboardID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = c.Header

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &result.Dashboard, nil
}

func (c *grafanaClient) GetPanels(dashboard models.Dashboard, uid string, params url.Values) map[int]models.Panel {
	panels := make(map[int]models.Panel)
	screen, err := strconv.ParseInt(params.Get("screen"), 0, 16)
	if err != nil {
		log.Printf("Failed to get the screen widht, defaulting to 1686")
		screen = 1686
	}
	params.Del("screen")
	p := params.Encode()

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
				Tag:     panel.Tag,
			}
			mu.Unlock()
			continue
		}
		wg.Add(1)
		go func(panel models.Panel) {
			defer wg.Done()
			url := getImageURL(c.BaseURL, uid, panel, p, screen)
			localImagePath := filepath.Join("../../static/images", filepath.Base(strconv.Itoa(panel.ID))) + ".png"
			err := c.downloadImage(url, localImagePath, panel.GridPos)
			if err != nil {
				log.Printf("Failed to download image for panel %d: %v", panel.ID, err)
				return
			}
			mu.Lock()
			panels[panel.ID] = models.Panel{
				ID:      panel.ID,
				URL:     "/static/images/" + strconv.Itoa(panel.ID) + ".png",
				GridPos: panel.GridPos,
				Tag:     panel.Tag,
			}
			mu.Unlock()
		}(panel)
	}
	wg.Wait()
	return panels
}

func getImageURL(url string, id string, panel models.Panel, params string, screen int64) string {
	width := getWidth(panel.GridPos.W, int(screen))
	height := getHeight(panel.GridPos.H)
	return fmt.Sprintf("%s/render/d-solo/%s/?panelId=%d&width=%d&height=%d&%s", url, id, panel.ID, width, height, params)
}

func (c *grafanaClient) downloadImage(url, filePath string, pos models.GridPos) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = c.Header

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

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	png.Encode(file, img)

	_, err = io.Copy(file, resp.Body)
	return err
}

func (c *grafanaClient) DeleteImages(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range files {
		if strings.HasSuffix(name, ".png") {
			err = os.Remove(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
