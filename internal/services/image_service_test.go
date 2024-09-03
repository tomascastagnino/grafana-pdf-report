package services

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients/mocks"
)

func TestImageService_FetchAndStoreImage(t *testing.T) {
	mockGrafanaClient := new(mocks.MockGrafanaClient)
	imageService := NewImageService(mockGrafanaClient)

	dashboardID := "test-dashboard"
	params := url.Values{}
	r := http.Request{}

	expectedPath := "/static/images/test_image.png"

	mockGrafanaClient.On("GetPanelImage", dashboardID, params, r.Header).Return(expectedPath, nil)

	path, err := imageService.FetchAndStoreImage(dashboardID, params, r)

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, path)

	mockGrafanaClient.AssertExpectations(t)
}

func TestImageService_FetchAndStoreImage_Error(t *testing.T) {
	mockGrafanaClient := new(mocks.MockGrafanaClient)
	imageService := NewImageService(mockGrafanaClient)

	dashboardID := "test-dashboard"
	params := url.Values{}
	r := http.Request{}

	expectedError := errors.New("failed to fetch image")

	mockGrafanaClient.On("GetPanelImage", dashboardID, params, r.Header).Return("", expectedError)

	path, err := imageService.FetchAndStoreImage(dashboardID, params, r)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, "", path)

	mockGrafanaClient.AssertExpectations(t)
}
