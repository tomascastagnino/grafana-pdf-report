package services

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/services/mocks"
)

func TestPanelService_GetPanelsWithImages_Success(t *testing.T) {
	mockImageService := new(mocks.MockImageService)
	panelService := NewPanelService(mockImageService)

	dashboard := &models.Dashboard{
		UID: "test-dashboard",
		Panels: []models.Panel{
			{ID: 1, GridPos: models.GridPos{W: 4, H: 3}},
			{ID: 2, GridPos: models.GridPos{W: 6, H: 3}},
		},
	}

	r := http.Request{}
	mockImageService.On("DeleteImages", internal.ImageDir).Return(nil)
	mockImageService.On("GetImagePath", dashboard.UID, dashboard.Panels[0], r).Return("/static/images/1.png", nil)
	mockImageService.On("GetImagePath", dashboard.UID, dashboard.Panels[1], r).Return("/static/images/2.png", nil)

	panels, err := panelService.GetPanelsWithImages(dashboard, r)
	assert.NoError(t, err)
	assert.Len(t, panels, 2)
	assert.ElementsMatch(t, []string{"/static/images/1.png", "/static/images/2.png"}, []string{panels[0].URL, panels[1].URL})

	mockImageService.AssertExpectations(t)
}

func TestPanelService_GetPanelsWithImages_Error(t *testing.T) {
	mockImageService := new(mocks.MockImageService)
	panelService := NewPanelService(mockImageService)

	dashboard := &models.Dashboard{
		UID: "test-dashboard",
		Panels: []models.Panel{
			{ID: 1, GridPos: models.GridPos{W: 4, H: 3}},
		},
	}

	r := http.Request{}

	mockImageService.On("DeleteImages", internal.ImageDir).Return(nil)
	mockImageService.On("GetImagePath", dashboard.UID, dashboard.Panels[0], r).Return("", errors.New("failed to get image"))

	panels, err := panelService.GetPanelsWithImages(dashboard, r)

	assert.Error(t, err)
	assert.Nil(t, panels)

	mockImageService.AssertExpectations(t)
}

func TestPanelService_GetPanel_Success(t *testing.T) {
	mockImageService := new(mocks.MockImageService)
	panelService := NewPanelService(mockImageService)

	dashboardID := "test-dashboard"
	panelID := "1"

	r := http.Request{
		URL: &url.URL{
			RawQuery: "w=500&h=400",
		},
	}

	panel := models.Panel{
		ID:      1,
		GridPos: models.GridPos{H: 400, W: 500, X: 0, Y: 0},
	}

	expectedImagePath := "/static/images/1.png"

	mockImageService.On("GetImagePath", dashboardID, panel, r).Return(expectedImagePath, nil)

	resultPanel, err := panelService.GetPanel(dashboardID, panelID, r)

	assert.NoError(t, err)
	assert.Equal(t, expectedImagePath, resultPanel.URL)

	mockImageService.AssertExpectations(t)
}

func TestPanelService_GetPanel_Error(t *testing.T) {
	mockImageService := new(mocks.MockImageService)
	panelService := NewPanelService(mockImageService)

	dashboardID := "test-dashboard"
	panelID := "1"
	r := http.Request{
		URL: &url.URL{
			RawQuery: "w=4&h=3",
		},
	}
	panel := models.Panel{
		ID:      1,
		GridPos: models.GridPos{W: 4, H: 3, X: 0, Y: 0},
	}

	expectedError := errors.New("failed to get image")

	mockImageService.On("GetImagePath", dashboardID, panel, r).Return("", expectedError)

	returnedPanel, err := panelService.GetPanel(dashboardID, panelID, r)

	assert.Contains(t, err.Error(), "failed to get image")
	assert.Equal(t, "", returnedPanel.URL)

	mockImageService.AssertExpectations(t)
}
