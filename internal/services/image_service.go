package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/utils"
)

type ImageServiceInterface interface {
	FetchAndStoreImage(dID string, params url.Values, r http.Request) (string, error)
	DeleteImages(dir string) error
	GetImagePath(dashboardID string, panel models.Panel, r http.Request) (string, error)
}

type ImageService struct {
	grafanaClient clients.GrafanaClient
}

func NewImageService(client clients.GrafanaClient) ImageServiceInterface {
	return &ImageService{grafanaClient: client}
}

func (s *ImageService) FetchAndStoreImage(dID string, params url.Values, r http.Request) (string, error) {
	// Fetch the image using the GrafanaClient and return the file path
	return s.grafanaClient.GetPanelImage(dID, params, r.Header)
}

func (s *ImageService) DeleteImages(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("failed to open directory %s: %w", dir, err)
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, name := range files {
		if strings.HasSuffix(name, ".png") {
			err = os.Remove(filepath.Join(dir, name))
			if err != nil {
				return fmt.Errorf("failed to remove image %s: %w", name, err)
			}
		}
	}
	return nil
}

func (s *ImageService) GetImagePath(dashboardID string, panel models.Panel, r http.Request) (string, error) {
	params := buildImageParams(panel, r.URL.Query())
	return s.FetchAndStoreImage(dashboardID, params, r)
}

func buildImageParams(panel models.Panel, params url.Values) url.Values {
	screen, _ := strconv.Atoi(params.Get("screen"))
	panelID := strconv.Itoa(panel.ID)
	width := strconv.Itoa(utils.GetWidth(panel.GridPos.W, int(screen)))
	height := strconv.Itoa(utils.GetHeight(panel.GridPos.H))

	params.Add("panelId", panelID)
	params.Add("width", width)
	params.Add("height", height)

	return params
}
