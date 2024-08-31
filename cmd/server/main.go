package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/clients"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/handlers"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/services"
)

func main() {
	os.MkdirAll(internal.ImageDir, os.ModePerm)

	grafanaClient := clients.NewGrafanaAPIClient(
		internal.BaseURL,
		&http.Client{},
		nil, // Optional headers
	)

	// Services
	imageService := services.NewImageService(grafanaClient)
	panelService := services.NewPanelService(imageService)
	dashboardService := services.NewDashboardService(grafanaClient, panelService)

	// Handlers
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	reportHandler := handlers.NewReportHandler()

	// Router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc(internal.DashboardPath, dashboardHandler.GetDashboard).Methods("GET")
	r.HandleFunc(internal.ReportPath, reportHandler.ServeReportPage).Methods("GET")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(internal.StaticDir))))
	r.PathPrefix("/node_modules/").Handler(http.StripPrefix("/node_modules/", http.FileServer(http.Dir(internal.NodeModulesDir))))

	// Start the server
	log.Printf("Server listening on port %s", ":9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}
