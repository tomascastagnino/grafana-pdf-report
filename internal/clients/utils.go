package clients

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/tomascastagnino/grafana-pdf-reporter/internal/models"
)

// getWidth calculates the width in pixels a panel's image should have.
// Each Grafana dashboard has a 16px padding left and a 16px padding right.
// Each Grafana panel has a 8px margin between each other. So, for a w:4 panel
// it's width is composed of the size of 4 panels + 3 * 8px (the 3 margins between
// each panel). The actual width of a w:1 panels is: the width of the window - 16 * 2
// (padding left and right) divided by 23 (24 is the max w:? value)
func getWidth(x int, screenWidth int) int {
	unit := (float32(screenWidth) - 32 - 23*8.33) / 24
	width := unit*float32(x) + (float32(x)-1)*8.33
	return int(width)
}

func getHeight(y int) int {
	return 30*y + 8*(y-1)
}

func imgPath(params url.Values) string {
	s := params.Get("panelId")
	w := params.Get("width")
	h := params.Get("height")
	path := fmt.Sprintf("%s_%s_%s", s, w, h)
	return filepath.Join("/static/images", filepath.Base(path)) + ".png"
}

func buildParams(r http.Request, panel models.Panel) url.Values {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	screen, _ := strconv.Atoi(params.Get("screen"))
	pID := strconv.Itoa(panel.ID)
	width := strconv.Itoa(getWidth(panel.GridPos.W, int(screen)))
	height := strconv.Itoa(getHeight(panel.GridPos.H))
	params.Add("panelId", pID)
	params.Add("width", width)
	params.Add("height", height)
	return params
}
