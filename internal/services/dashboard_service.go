package services

import (
	"net/http"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type DashboardService struct {
	grafanaClient clients.GrafanaClient
	panelService  *PanelService
}

func NewDashboardService(client clients.GrafanaClient, panelService *PanelService) *DashboardService {
	return &DashboardService{grafanaClient: client, panelService: panelService}
}

func (s *DashboardService) GetDashboard(dashboardID string, r *http.Request) (*models.Dashboard, error) {
	// Fetch dashboard metadata
	dashboard, err := s.grafanaClient.GetDashboard(dashboardID, r.Header)
	if err != nil {
		return nil, err
	}

	// Get panels with images using the PanelService
	panels, err := s.panelService.GetPanelsWithImages(dashboard, *r)
	if err != nil {
		return nil, err
	}

	// Assign panels to the dashboard
	dashboard.Panels = panels

	return dashboard, nil
}

func (s *DashboardService) ListDashboards() ([]models.Dashboard, error) {
	panic("TODO: To be implemented in the future.")
}