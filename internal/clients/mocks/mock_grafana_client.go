package mocks

import (
	"net/http"
	"net/url"

	"github.com/stretchr/testify/mock"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type MockGrafanaClient struct {
	mock.Mock
}

func (m *MockGrafanaClient) GetDashboard(dashboardID string, h http.Header) (*models.Dashboard, error) {
	args := m.Called(dashboardID, h)
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *MockGrafanaClient) GetPanelImage(dID string, params url.Values, h http.Header) (string, error) {
	args := m.Called(dID, params, h)
	return args.String(0), args.Error(1)
}

func (m *MockGrafanaClient) GetAllDashboards(h http.Header) ([]models.Dashboard, error) {
	args := m.Called(h)
	return args.Get(0).([]models.Dashboard), args.Error(1)
}
