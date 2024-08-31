package internal

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	APIVersion       string = "/api/v1/"
	ReportPath       string = APIVersion + "report/{id}/"
	DashboardPath    string = APIVersion + "dashboard/{id}/"
	RefreshPanelPath string = ReportPath + "refresh_panel/"
	GrafanaURL       string
	BaseURL          string
	BaseDir          string
	StaticDir        string
	ImageDir         string
	WebImageDir      string
	NodeModulesDir   string
	DashboardURL     string
	ImageRendererURL string
	ConfigFilePath   string
	ChannelNum       int
)

func init() {
	cfg, err := ini.Load(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	BaseDir = cfg.Section("paths").Key("BaseDir").MustString("/app")

	StaticDir = filepath.Join(BaseDir, "static")
	ImageDir = filepath.Join(StaticDir, "images")
	WebImageDir = "/static/images"
	NodeModulesDir = filepath.Join(BaseDir, "node_modules")
	ConfigFilePath = filepath.Join(BaseDir, "config.ini")

	GrafanaURL = cfg.Section("server").Key("GrafanaURL").MustString("http://grafana:3000")
	DashboardURL = cfg.Section("url").Key("DashboardURL").MustString("/api/dashboards/uid/")
	BaseURL = GrafanaURL + DashboardURL
	ImageRendererURL = GrafanaURL + cfg.Section("url").Key("ImageRendererURL").MustString("/render/d-solo")

	ChannelNum = cfg.Section("channels").Key("ChannelNum").MustInt(10)
}
