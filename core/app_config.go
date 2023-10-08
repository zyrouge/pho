package core

import (
	"fmt"
	"path"

	"github.com/zyrouge/pho/utils"
)

type AppConfig struct {
	Id      string   `json:"Id"`
	Name    string   `json:"Name"`
	Version string   `json:"Version"`
	Source  SourceId `json:"Source"`
	Paths   AppPaths `json:"Paths"`
}

type SourceId string

type AppPaths struct {
	Dir          string
	Config       string
	SourceConfig string
	AppImage     string
	Icon         string
	Desktop      string
}

func ReadAppConfig(configPath string) (*AppConfig, error) {
	appConfig, err := utils.ReadJsonFile[AppConfig](configPath)
	if err != nil {
		return nil, err
	}
	// TODO: remove this in future versions
	// patch this temporarily to not break on user
	if appConfig.Paths.AppImage == "" {
		config, err := ReadConfig()
		if err != nil {
			return nil, err
		}
		appConfig.Paths = *ConstructAppPaths(config, appConfig.Id, appConfig.Name)
	}
	return appConfig, nil
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

func ConstructAppName(appName string) string {
	return utils.CleanName(appName)
}
func ConstructAppPaths(config *Config, appId string, appName string) *AppPaths {
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
