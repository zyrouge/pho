package core

import (
	"errors"

	"github.com/zyrouge/pho/utils"
)

const HttpSourceId SourceId = "http"

type HttpSource struct{}

func ReadHttpSourceConfig(configPath string) (*HttpSource, error) {
	return utils.ReadJsonFile[HttpSource](configPath)
}

func (*HttpSource) SupportsUpdates() bool {
	return false
}

func (*HttpSource) CheckUpdate(app *AppConfig) (bool, *SourceUpdate, error) {
	return false, nil, errors.New("http source does not support updates")
}
