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

	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
	"gopkg.in/ini.v1"
)

var (
	GrafanaURL   string
	dashboardURL = internal.GDashboardURL
)

type grafanaClient struct {
	http.Client

	BaseURL      string
	dashboardURL string
	Header       http.Header
	ChannelNum   int
}

func GetGrafanaClient(r *http.Header) *grafanaClient {
	cfg, err := ini.Load(internal.ConfigFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	num, err := cfg.Section("channels").Key("ChannelNum").Int()
	if err != nil || num <= 0 {
		num = 10 // default number of channels
	}

	return &grafanaClient{
		BaseURL:      getGrafanaURL(),
		dashboardURL: dashboardURL,
		Header:       r.Clone(),
		ChannelNum:   num,
	}
}

func getGrafanaURL() string {
	cfg, err := ini.Load(internal.ConfigFilePath)
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.Header

	resp, err := c.Do(req)
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
		return nil, fmt.Errorf("failed to decode dashboard response: %w", err)
	}

	return &result.Dashboard, nil
}

func (c *grafanaClient) GetPanels(dashboard models.Dashboard, r http.Request) (map[int]models.Panel, error) {
	panels := make(map[int]models.Panel)
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, c.ChannelNum)
	errorChannel := make(chan error, len(dashboard.Panels))
	defer close(errorChannel)

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
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			path, err := c.downloadImage(panel, r)
			if err != nil {
				errorChannel <- fmt.Errorf("failed to download image for panel %d: %w", panel.ID, err)
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

	if len(errorChannel) > 0 {
		var errs error
		for err := range errorChannel {
			if errs == nil {
				errs = err
			} else {
				errs = fmt.Errorf("%v; %w", errs, err)
			}
		}
		return nil, errs
	}

	return panels, nil
}

func (c *grafanaClient) getImageURL(p url.Values) string {
	return fmt.Sprintf("%s/%s/%s/?%s", c.BaseURL, internal.ImageRendererURL, p.Get("dashboardId"), p.Encode())
}

func (c *grafanaClient) downloadImage(panel models.Panel, r http.Request) (string, error) {
	base := "../.."
	p := buildParams(r, panel)
	url := c.getImageURL(p)
	path := filepath.Join(base, imgPath(p))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for image download: %w", err)
	}
	req.Header = c.Header

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
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
	png.Encode(file, img)

	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("failed to encode image to PNG: %w", err)
	}

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
		return "", fmt.Errorf("failed to download image for panel %d: %w", panel.ID, err)
	}
	return path, nil
}
