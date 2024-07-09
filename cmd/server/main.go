package main

import (
    "log"
    "net/http"

	"github.com/tomascastagnino/grafana-pdf-report/internal/handlers"
)

const apiVersion = "/api/v1/"

func main() {
    http.HandleFunc(apiVersion+"report/", handlers.HandleReport)
    http.HandleFunc(apiVersion+"report/data/", handlers.HandleReportData)
    http.HandleFunc("/generate-pdf", handlers.HandleGeneratePDF)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))
    log.Printf("Server listening on port %s", ":9090")
    log.Fatal(http.ListenAndServe(":9090", nil))
}