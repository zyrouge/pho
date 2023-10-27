package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type GithubApiRelease struct {
	ApiUrl     string                  `json:"url"`
	HtmlUrl    string                  `json:"html_url"`
	TagName    string                  `json:"tag_name"`
	Draft      bool                    `json:"draft"`
	PreRelease bool                    `json:"prerelease"`
	Assets     []GithubApiReleaseAsset `json:"assets"`
}

type GithubApiReleaseAsset struct {
	ApiUrl      string `json:"url"`
	DownloadUrl string `json:"browser_download_url"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
}

func GithubApiFetchReleases(username string, reponame string) (*[]GithubApiRelease, error) {
	return RequestGithubApi[[]GithubApiRelease](
		"GET",
		fmt.Sprintf("/repos/%s/%s/releases", username, reponame),
	)
}

func GithubApiFetchLatestPreRelease(username string, reponame string) (*GithubApiRelease, error) {
	releases, err := GithubApiFetchReleases(username, reponame)
	if err != nil {
		return nil, err
	}
	for _, x := range *releases {
		if x.PreRelease {
			return &x, nil
		}
	}
	return nil, errors.New("no prerelease found")
}

func GithubApiFetchLatestAny(username string, reponame string) (*GithubApiRelease, error) {
	releases, err := GithubApiFetchReleases(username, reponame)
	if err != nil {
		return nil, err
	}
	for _, x := range *releases {
		if !x.Draft {
			return &x, nil
		}
	}
	return nil, errors.New("no non-draft releases found")
}

func GithubApiFetchLatestRelease(username string, reponame string) (*GithubApiRelease, error) {
	return RequestGithubApi[GithubApiRelease](
		"GET",
		fmt.Sprintf("/repos/%s/%s/releases/latest", username, reponame),
	)
}

func GithubApiFetchTaggedRelease(username string, reponame string, tag string) (*GithubApiRelease, error) {
	return RequestGithubApi[GithubApiRelease](
		"GET",
		fmt.Sprintf("/repos/%s/%s/releases/tags/%s", username, reponame, tag),
	)
}

func (asset *GithubApiReleaseAsset) ToAsset() *Asset {
	return &Asset{
		Source:   asset.DownloadUrl,
		Size:     asset.Size,
		Download: NetworkAssetDownload(asset.DownloadUrl),
	}
}

var GithubRepoUrlRegex = regexp.MustCompile(`^([^\/]+)\/([^\/]+)$`)

func ParseGithubRepoUrl(url string) (bool, string, string) {
	url = strings.TrimPrefix(url, "https://github.com/")
	matches := GithubRepoUrlRegex.FindStringSubmatch(url)
	if matches == nil {
		return false, "", ""
	}
	return true, matches[1], matches[2]
}

func RequestGithubApi[T any](method string, route string) (*T, error) {
	url := fmt.Sprintf("https://api.github.com%s", route)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(
			"github api response returned status %d with message \"%s\"",
			res.StatusCode,
			res.Status,
		)
	}
	output := new(T)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
