package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type PanelService struct {
	imageService *ImageService
}

func NewPanelService(imageService *ImageService) *PanelService {
	return &PanelService{imageService: imageService}
}

func (s *PanelService) GetPanelsWithImages(dashboard *models.Dashboard, r http.Request) ([]models.Panel, error) {
	s.imageService.DeleteImages(internal.ImageDir)

	var panels []models.Panel
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, internal.ChannelNum)
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

			imagePath, err := s.imageService.GetImagePath(dashboard.UID, panel, r)
			if err != nil {
				errorChannel <- fmt.Errorf("failed to download image for panel %d: %w", panel.ID, err)
				return
			}
			mu.Lock()
			panel.URL = imagePath
			panels = append(panels, panel)
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

func (s *PanelService) GetPanel(dashboardID string, panelID string, r http.Request) (models.Panel, error) {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	w, _ := strconv.Atoi(params.Get("w"))
	h, _ := strconv.Atoi(params.Get("h"))
	id, _ := strconv.Atoi(panelID)

	pos := models.GridPos{
		H: h,
		W: w,
		X: 0,
		Y: 0,
	}
	panel := models.Panel{
		ID:      id,
		GridPos: pos,
	}

	imagePath, err := s.imageService.GetImagePath(dashboardID, panel, r)
	if err != nil {
		return panel, fmt.Errorf("failed to get image path: %w", err)
	}

	panel.URL = imagePath

	return panel, nil
}
