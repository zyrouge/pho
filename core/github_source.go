package core

import (
	"errors"
	"fmt"

	"github.com/zyrouge/pho/utils"
)

const GithubSourceId SourceId = "github"

type GithubSourceRelease string

const (
	GithubSourceReleaseLatest     GithubSourceRelease = "latest"
	GithubSourceReleasePreRelease GithubSourceRelease = "prerelease"
	GithubSourceReleaseAny        GithubSourceRelease = "any"
)

type GithubSource struct {
	UserName string              `json:"UserName"`
	RepoName string              `json:"RepoName"`
	Release  GithubSourceRelease `json:"Release"`
}

func ReadGithubSourceConfig(configPath string) (*GithubSource, error) {
	return utils.ReadJsonFile[GithubSource](configPath)
}

func (source *GithubSource) FetchAptRelease() (*GithubApiRelease, error) {
	return source.FetchAptLatestRelease()
}

func (source *GithubSource) FetchAptLatestRelease() (*GithubApiRelease, error) {
	switch source.Release {
	case GithubSourceReleaseLatest:
		return GithubApiFetchLatestRelease(source.UserName, source.RepoName)

	case GithubSourceReleasePreRelease:
		return GithubApiFetchLatestPreRelease(source.UserName, source.RepoName)

	case GithubSourceReleaseAny:
		return GithubApiFetchLatestAny(source.UserName, source.RepoName)

	default:
		return nil, errors.New("invalid github source release type")
	}
}

func (release *GithubApiRelease) ChooseAptAsset() (AppImageAssetMatch, *GithubApiReleaseAsset) {
	return ChooseAptAppImageAsset(
		release.Assets,
		func(x *GithubApiReleaseAsset) string {
			return x.Name
		},
	)
}

func (source *GithubSource) SupportUpdates() bool {
	return true
}

func (source *GithubSource) CheckUpdate(app *AppConfig, reinstall bool) (*SourceUpdate, error) {
	release, err := source.FetchAptLatestRelease()
	if err != nil {
		return nil, err
	}
	if app.Version == release.TagName && !reinstall {
		return nil, nil
	}
	matchScore, asset := release.ChooseAptAsset()
	if matchScore == AppImageAssetNoMatch {
		return nil, fmt.Errorf("no valid asset in github tag %s", release.TagName)
	}
	update := &SourceUpdate{
		Version:    release.TagName,
		MatchScore: matchScore,
		Asset:      asset.ToAsset(),
	}
	return update, nil
}
