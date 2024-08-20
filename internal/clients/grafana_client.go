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

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
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

func (c *grafanaClient) GetPanels(dashboard models.Dashboard, r http.Request) map[int]models.Panel {
	panels := make(map[int]models.Panel)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, panel := range dashboard.Panels {
		if panel.Tag == "remove" {
			continue
		}

		// I need to sanitize the HTML
		// if panel.Type == "text" {
		// 	mu.Lock()
		// 	panels[panel.ID] = models.Panel{
		// 		ID:      panel.ID,
		// 		Type:    panel.Type,
		// 		GridPos: panel.GridPos,
		// 		Options: panel.Options,
		// 		Tag:     panel.Tag,
		// 	}
		// 	mu.Unlock()
		// 	continue
		// }

		wg.Add(1)
		go func(panel models.Panel) {
			defer wg.Done()
			path, err := c.downloadImage(panel, r)
			if err != nil {
				log.Printf("Failed to download image for panel %d: %v", panel.ID, err)
				return
			}
			mu.Lock()
			panels[panel.ID] = models.Panel{
				ID:      panel.ID,
				URL:     path,
				GridPos: panel.GridPos,
				Tag:     panel.Tag,
			}
			mu.Unlock()
		}(panel)
	}
	wg.Wait()
	return panels
}

func (c *grafanaClient) getImageURL(p url.Values) string {
	return fmt.Sprintf("%s/render/d-solo/%s/?%s", c.BaseURL, p.Get("dashboardId"), p.Encode())
}

func (c *grafanaClient) downloadImage(panel models.Panel, r http.Request) (string, error) {
	base := "../.."
	p := buildParams(r, panel)
	url := c.getImageURL(p)
	log.Println(url)
	path := filepath.Join(base, imgPath(p))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header = c.Header

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return "", err
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	png.Encode(file, img)

	_, err = io.Copy(file, resp.Body)
	return imgPath(p), err
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

func (c *grafanaClient) GetRefreshedPanelURL(r http.Request) (string, error) {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	w, _ := strconv.Atoi(params.Get("w"))
	h, _ := strconv.Atoi(params.Get("h"))
	id, _ := strconv.Atoi(params.Get("panelId"))
	pos := models.GridPos{
		H: h,
		W: w,
		X: 0,
		Y: 0,
	}
	panel := models.Panel{
		ID:      id,
		URL:     imgPath(params),
		GridPos: pos,
	}
	path, err := c.downloadImage(panel, r)
	if err != nil {
		log.Printf("Failed to download image for panel %d: %v", panel.ID, err)
		return "", err
	}
	return path, nil
}
