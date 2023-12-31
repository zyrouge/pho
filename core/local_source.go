package core

import (
	"errors"

	"github.com/zyrouge/pho/utils"
)

const LocalSourceId SourceId = "local"

type LocalSource struct{}

func ReadLocalSourceConfig(configPath string) (*LocalSource, error) {
	return utils.ReadJsonFile[LocalSource](configPath)
}

func (*LocalSource) SupportUpdates() bool {
	return false
}

func (*LocalSource) CheckUpdate(app *AppConfig, reinstall bool) (*SourceUpdate, error) {
	return nil, errors.New("local source does not support updates")
}
