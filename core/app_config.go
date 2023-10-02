package core

import (
	"fmt"
	"path"

	"github.com/zyrouge/pho/utils"
)

type AppConfig struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	AppImage string `json:"AppImage"`
	Icon     string `json:"Icon"`
	Version  string `json:"Version"`
	Source   string `json:"Source"`
}

func ReadAppConfig(configPath string) (*AppConfig, error) {
	return utils.ReadJsonFile[AppConfig](configPath)
}

func SaveAppConfig(configPath string, config *AppConfig) error {
	return utils.WriteJsonFile[AppConfig](configPath, config)
}

func SaveAppSourceConfig[T any](configPath string, config T) error {
	return utils.WriteJsonFile[T](configPath, &config)
}

func ConstructAppId(owner string, appName string) string {
	raw := fmt.Sprintf("%s-%s", owner, appName)
	return utils.CleanId(raw)
}

type AppPaths struct {
	Dir          string
	Config       string
	SourceConfig string
	AppImage     string
	Icon         string
	Desktop      string
}

func GetAppPaths(config *Config, appId string, appName string) *AppPaths {
	appDir := path.Join(config.AppsDir, appId)
	cleanAppName := utils.CleanText(appName)
	if cleanAppName == "" {
		cleanAppName = appId
	}
	return &AppPaths{
		Dir:          appDir,
		Config:       path.Join(appDir, "config.zap.json"),
		SourceConfig: path.Join(appDir, "source.config.zap.json"),
		AppImage:     path.Join(appDir, fmt.Sprintf("%s.AppImage", cleanAppName)),
		Icon:         path.Join(appDir, fmt.Sprintf("%s.png", cleanAppName)),
		Desktop:      path.Join(config.DesktopDir, fmt.Sprintf("%s.desktop", appId)),
	}
}
