package core

import "errors"

type Source interface {
	SupportUpdates() bool
	CheckUpdate(app *AppConfig) (bool, *SourceUpdate, error)
}

type SourceUpdate struct {
	Version    string
	MatchScore AppImageAssetMatch
	*Asset
}

func ReadSourceConfig(sourceId SourceId, sourcePath string) (any, error) {
	switch sourceId {
	case GithubSourceId:
		return ReadGithubSourceConfig(sourcePath)

	case HttpSourceId:
		return ReadHttpSourceConfig(sourcePath)

	case LocalSourceId:
		return ReadLocalSourceConfig(sourcePath)

	default:
		return nil, errors.New("invalid source id")
	}
}

func CastSourceConfigAsSource(config any) (Source, error) {
	source, ok := config.(Source)
	if !ok {
		return nil, errors.New("config does not implement source")
	}
	return source, nil
}
