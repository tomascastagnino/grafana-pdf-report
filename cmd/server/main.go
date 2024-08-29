package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
	"github.com/tomascastagnino/grafana-pdf-reporter/internal/handlers"
)

func main() {
	os.Mkdir(internal.ImageDir, os.ModePerm)

	http.HandleFunc(internal.ReportPath, handlers.HandleReport)
	http.HandleFunc(internal.ReportDataPath, handlers.HandleReportData)
	http.HandleFunc(internal.RefreshPanelPath, handlers.HandleRefresh)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(internal.StaticDir))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir(internal.NodeModulesDir))))

	log.Printf("Server listening on port %s", ":9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
