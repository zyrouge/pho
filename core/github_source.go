package core

const GithubSourceId = "github"

type GithubSource struct {
	UserName   string `json:"UserName"`
	RepoName   string `json:"RepoName"`
	PreRelease bool   `json:"PreRelease"`
	TagName    string `json:"TagName"`
}

func (source *GithubSource) FetchAptRelease() (*GithubApiRelease, error) {
	if source.TagName != "" {
		return GithubApiFetchTaggedRelease(source.UserName, source.RepoName, source.TagName)
	}
	if source.PreRelease {
		return GithubApiFetchLatestPreRelease(source.UserName, source.RepoName)
	}
	return GithubApiFetchLatestRelease(source.UserName, source.RepoName)
}
