package core

import (
	"fmt"
	"path"

	"github.com/zyrouge/pho/utils"
)

type AppConfig struct {
	Id      string   `json:"Id"`
	Version string   `json:"Version"`
	Source  SourceId `json:"Source"`
	Paths   AppPaths `json:"Paths"`
}

type SourceId string

type AppPaths struct {
	Dir          string `json:"Dir"`
	Config       string `json:"Config"`
	SourceConfig string `json:"SourceConfig"`
	AppImage     string `json:"AppImage"`
	Icon         string `json:"Icon"`
	Desktop      string `json:"Desktop"`
}

func ReadAppConfig(configPath string) (*AppConfig, error) {
	return utils.ReadJsonFile[AppConfig](configPath)
}

func SaveAppConfig(configPath string, config *AppConfig) error {
	return utils.WriteJsonFile[AppConfig](configPath, config)
}

func SaveSourceConfig[T any](configPath string, config T) error {
	return utils.WriteJsonFile[T](configPath, &config)
}

func ConstructAppId(appName string) string {
	return utils.CleanId(appName)
}

func ConstructAppPaths(config *Config, appId string) *AppPaths {
	appDir := path.Join(config.AppsDir, appId)
	return &AppPaths{
		Dir:          appDir,
		Config:       path.Join(appDir, "config.pho.json"),
		SourceConfig: path.Join(appDir, "source.config.pho.json"),
		AppImage:     path.Join(appDir, fmt.Sprintf("%s.AppImage", appId)),
		Icon:         path.Join(appDir, fmt.Sprintf("%s.png", appId)),
		Desktop:      path.Join(config.DesktopDir, fmt.Sprintf("%s.desktop", appId)),
	}
}

func ConstructAppConfigPath(config *Config, appId string) string {
	appPaths := ConstructAppPaths(config, appId)
	return appPaths.Config
}
