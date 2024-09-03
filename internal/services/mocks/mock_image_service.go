package mocks

import (
	"net/http"
	"net/url"

	"github.com/stretchr/testify/mock"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) FetchAndStoreImage(dID string, params url.Values, r http.Request) (string, error) {
	args := m.Called(dID, params, r)
	return args.String(0), args.Error(1)
}

func (m *MockImageService) DeleteImages(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func (m *MockImageService) GetImagePath(dashboardID string, panel models.Panel, r http.Request) (string, error) {
	args := m.Called(dashboardID, panel, r)
	return args.String(0), args.Error(1)
}
