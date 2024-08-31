package services

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients"
)

type ImageService struct {
	grafanaClient clients.GrafanaClient
}

func NewImageService(client clients.GrafanaClient) *ImageService {
	return &ImageService{grafanaClient: client}
}

func (s *ImageService) FetchAndStoreImage(dID string, params url.Values, r http.Request) (string, error) {
	// Fetch the image using the GrafanaClient and return the file path
	return s.grafanaClient.GetPanelImage(dID, params, r.Header)
}

func (s *ImageService) DeleteImages(dir string) error {
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
