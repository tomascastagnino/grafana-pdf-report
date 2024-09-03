package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type MockPanelService struct {
	mock.Mock
}

func (m *MockPanelService) GetPanelsWithImages(dashboard *models.Dashboard, r http.Request) ([]models.Panel, error) {
	args := m.Called(dashboard, r)
	return args.Get(0).([]models.Panel), args.Error(1)
}

func (m *MockPanelService) GetPanel(dashboardID string, panelID string, r http.Request) (models.Panel, error) {
	args := m.Called(dashboardID, panelID, r)
	return args.Get(0).(models.Panel), args.Error(1)
}
