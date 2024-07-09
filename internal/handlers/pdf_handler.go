package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"

	"github.com/tomascastagnino/grafana-pdf-report/internal/clients"
    "github.com/jung-kurt/gofpdf"
)

func HandleGeneratePDF(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        Panels []string `json:"panels"`
    }

    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetFont("Arial", "", 14)

    client := clients.GetGrafanaClient()

    for _, panelURL := range requestData.Panels {
        req, err := http.NewRequest("GET", panelURL, nil)
        if err != nil {
            log.Printf("failed to create request: %v", err)
            http.Error(w, fmt.Sprintf("failed to create request: %v", err), http.StatusInternalServerError)
            return
        }

        req.Header.Set("Authorization", "Bearer "+client.ApiKey)

        resp, err := client.Do(req)
        if err != nil {
            log.Printf("failed to fetch image: %v", err)
            http.Error(w, fmt.Sprintf("failed to fetch image: %v", err), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            log.Printf("non-200 response: %d", resp.StatusCode)
            bodyBytes, _ := ioutil.ReadAll(resp.Body)
            log.Printf("response body: %s", bodyBytes)
            http.Error(w, fmt.Sprintf("non-200 response: %d", resp.StatusCode), http.StatusInternalServerError)
            return
        }

        contentType := resp.Header.Get("Content-Type")
        if !strings.Contains(contentType, "image/png") {
            bodyBytes, _ := ioutil.ReadAll(resp.Body)
            log.Printf("unexpected content type: %s, body: %s", contentType, string(bodyBytes))
            http.Error(w, fmt.Sprintf("unexpected content type: %s", contentType), http.StatusInternalServerError)
            return
        }

        img, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("failed to read image: %v", err)
            http.Error(w, fmt.Sprintf("failed to read image: %v", err), http.StatusInternalServerError)
            return
        }

        options := gofpdf.ImageOptions{
            ImageType:             "PNG",
            ReadDpi:               true,
            AllowNegativePosition: false,
        }
        pdf.AddPage()
        pdf.RegisterImageOptionsReader(panelURL, options, bytes.NewReader(img))
        pdf.ImageOptions(panelURL, 10, 10, 190, 0, false, options, 0, "")
    }

    var buf bytes.Buffer
    err = pdf.Output(&buf)
    if err != nil {
        log.Printf("failed to generate PDF: %v", err)
        http.Error(w, fmt.Sprintf("failed to generate PDF: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "attachment; filename=selected_panels.pdf")
    w.Write(buf.Bytes())
}