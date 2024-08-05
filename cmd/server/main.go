package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tomascastagnino/grafana-pdf-report/internal/handlers"
)

const apiVersion = "/api/v1/"

func main() {
	os.Mkdir("../../static/images", os.ModePerm)

	http.HandleFunc(apiVersion+"report/", handlers.HandleReport)
	http.HandleFunc(apiVersion+"report/data/", handlers.HandleReportData)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("../../node_modules"))))
	log.Printf("Server listening on port %s", ":9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
