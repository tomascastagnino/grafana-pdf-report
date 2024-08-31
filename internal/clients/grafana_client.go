package clients

import (
	"net/http"
	"net/url"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

type GrafanaClient interface {
	GetDashboard(dashboardID string, h http.Header) (*models.Dashboard, error)
	GetPanelImage(dID string, params url.Values, h http.Header) (string, error)
	GetAllDashboards() ([]models.Dashboard, error)
}
