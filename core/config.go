package core

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/zyrouge/pho/utils"
)

type Config struct {
	AppsDir    string            `json:"AppsDir"`
	DesktopDir string            `json:"DesktopDir"`
	Installed  map[string]string `json:"Installed"`
}

var cachedConfig *Config

func GetConfigPath() (string, error) {
	xdgConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := path.Join(xdgConfigDir, AppCodeName, "config.json")
	return configPath, nil
}

func ReadConfig() (*Config, error) {
	cachedConfig = nil
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := utils.ReadJsonFile[Config](configPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf(
			"config file does not exist, use %s %s to initiate the setup",
			color.CyanString(AppExecutableName),
			color.CyanString("init"),
		)
	}
	if err != nil {
		return nil, err
	}
	cachedConfig = config
	return config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	err = utils.WriteJsonFileAtomic[Config](configPath, config)
	if err != nil {
		return err
	}
	cachedConfig = config
	return nil
}

func GetConfig() (*Config, error) {
	if cachedConfig == nil {
		return ReadConfig()
	}
	return cachedConfig, nil
}
