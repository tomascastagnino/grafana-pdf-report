package services

import (
	"fmt"
	"net/http"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type DashboardServiceInterface interface {
	GetDashboard(dashboardID string, r *http.Request) (*models.Dashboard, error)
	ListDashboards(r *http.Request) ([]models.Dashboard, error)
}

type DashboardService struct {
	grafanaClient clients.GrafanaClient
	panelService  PanelServiceInterface
}

func NewDashboardService(client clients.GrafanaClient, panelService PanelServiceInterface) *DashboardService {
	return &DashboardService{grafanaClient: client, panelService: panelService}
}

func (s *DashboardService) GetDashboard(dashboardID string, r *http.Request) (*models.Dashboard, error) {
	// Fetch dashboard metadata
	dashboard, err := s.grafanaClient.GetDashboard(dashboardID, r.Header)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dashboard %s: %w", dashboardID, err)
	}

	// Get panels with images using the PanelService
	panels, err := s.panelService.GetPanelsWithImages(dashboard, *r)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch panels for dashboard %s: %w", dashboardID, err)
	}

	// Assign panels to the dashboard
	dashboard.Panels = panels

	return dashboard, nil
}

func (s *DashboardService) ListDashboards(r *http.Request) ([]models.Dashboard, error) {
	// Get all the dashboards
	list, err := s.grafanaClient.GetAllDashboards(r.Header)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all the dashboards: %w", err)
	}
	return list, nil
}
