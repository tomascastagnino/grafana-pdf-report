package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

func (h *ReportHandler) ServeReportPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(internal.StaticDir, "index.html"))
}
