package handlers

import (
    "encoding/json"
    "net/http"
	"path/filepath"
    "os"
	"log"
	"strings"

	"github.com/tomascastagnino/grafana-pdf-report/internal/models"
    "github.com/jung-kurt/gofpdf"
)

func HandleGeneratePDF(w http.ResponseWriter, r *http.Request) {
	var requestData models.PdfRequest

    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()

    for _, panel := range requestData.Panels {
        x := float64(panel.GridPos.X) * 10
        y := float64(panel.GridPos.Y) * 10
        w := float64(panel.GridPos.W) * 10
        h := float64(panel.GridPos.H) * 10

        if panel.URL != "" {
			url := "../.." + panel.URL
            pdf.Image(url, x, y, w, h, false, "", 0, "")
        }
    }

    w.Header().Set("Content-Type", "application/pdf")
    err = pdf.Output(w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	err = DeleteImages("../../static/images")
    if err != nil {
        log.Printf("Error deleting images: %v", err)
    }
}

func DeleteImages(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()

    files, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }

    for _, name := range files {
        if strings.HasSuffix(name, ".png") {
            err = os.Remove(filepath.Join(dir, name))
            if err != nil {
                return err
            }
        }
    }
    return nil
}