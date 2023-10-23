package core

import (
	"fmt"
	"path"
	"strings"

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
	appConfig, err := utils.ReadJsonFile[AppConfig](configPath)
	if err != nil {
		return nil, err
	}
	// TODO: remove this in future versions
	// patched this temporarily to not break on user
	if appConfig.Paths.AppImage == "" {
		config, err := ReadConfig()
		if err != nil {
			return nil, err
		}
		appConfig.Paths = *ConstructAppPaths(config, appConfig.Id)
		appConfig.Paths.Config = strings.Replace(
			appConfig.Paths.Config,
			".pho.json",
			".zap.json",
			1,
		)
		appConfig.Paths.SourceConfig = strings.Replace(
			appConfig.Paths.SourceConfig,
			".pho.json",
			".zap.json",
			1,
		)
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
