package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	clientsMocks "github.com/tomascastagnino/grafana-pdf-reporter/internal/clients/mocks"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
	servicesMocks "github.com/tomascastagnino/grafana-pdf-reporter/internal/services/mocks"
)

func TestDashboardService_GetDashboard_Success(t *testing.T) {
	mockPanelService := new(servicesMocks.MockPanelService)
	mockGrafanaClient := new(clientsMocks.MockGrafanaClient)

	dashboardService := NewDashboardService(mockGrafanaClient, mockPanelService)

	dashboardID := "test-dashboard"
	r := http.Request{}
	expectedDashboard := &models.Dashboard{
		UID:    dashboardID,
		Panels: []models.Panel{},
	}

	mockGrafanaClient.On("GetDashboard", dashboardID, r.Header).Return(expectedDashboard, nil)
	mockPanelService.On("GetPanelsWithImages", expectedDashboard, r).Return(expectedDashboard.Panels, nil)

	dashboard, err := dashboardService.GetDashboard(dashboardID, &r)

	assert.NoError(t, err)
	assert.Equal(t, expectedDashboard, dashboard)

	mockGrafanaClient.AssertExpectations(t)
	mockPanelService.AssertExpectations(t)
}

func TestDashboardService_GetDashboard_ErrorFetchingDashboard(t *testing.T) {
	mockPanelService := new(servicesMocks.MockPanelService)
	mockGrafanaClient := new(clientsMocks.MockGrafanaClient)

	dashboardService := NewDashboardService(mockGrafanaClient, mockPanelService)

	dashboardID := "test-dashboard"
	r := &http.Request{}
	expectedError := errors.New("failed to fetch dashboard")

	mockGrafanaClient.On("GetDashboard", dashboardID, r.Header).Return((*models.Dashboard)(nil), expectedError)

	dashboard, err := dashboardService.GetDashboard(dashboardID, r)

	assert.Error(t, err)
	assert.Nil(t, dashboard)
	assert.Contains(t, err.Error(), "failed to fetch dashboard")

	mockGrafanaClient.AssertExpectations(t)
}

func TestDashboardService_ListDashboards_Success(t *testing.T) {
	mockPanelService := new(servicesMocks.MockPanelService)
	mockGrafanaClient := new(clientsMocks.MockGrafanaClient)

	dashboardService := NewDashboardService(mockGrafanaClient, mockPanelService)

	r := http.Request{}
	expectedDashboards := []models.Dashboard{
		{UID: "dashboard1", Title: "Dashboard 1"},
		{UID: "dashboard2", Title: "Dashboard 2"},
	}

	mockGrafanaClient.On("GetAllDashboards", r.Header).Return(expectedDashboards, nil)

	dashboards, err := dashboardService.ListDashboards(&r)

	assert.NoError(t, err)
	assert.Equal(t, expectedDashboards, dashboards)

	mockGrafanaClient.AssertExpectations(t)
}

func TestDashboardService_ListDashboards_Error(t *testing.T) {
	mockPanelService := new(servicesMocks.MockPanelService)
	mockGrafanaClient := new(clientsMocks.MockGrafanaClient)

	dashboardService := NewDashboardService(mockGrafanaClient, mockPanelService)

	r := &http.Request{}
	expectedError := errors.New("failed to fetch all the dashboards")

	mockGrafanaClient.On("GetAllDashboards", r.Header).Return(([]models.Dashboard)(nil), expectedError)

	dashboards, err := dashboardService.ListDashboards(r)

	assert.Error(t, err)
	assert.Nil(t, dashboards)
	assert.Contains(t, err.Error(), "failed to fetch all the dashboards")

	mockGrafanaClient.AssertExpectations(t)
}
