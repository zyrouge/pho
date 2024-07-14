package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var githubSourceReleaseStrings = []string{
	string(core.GithubSourceReleaseLatest),
	string(core.GithubSourceReleasePreRelease),
	string(core.GithubSourceReleaseTagged),
	string(core.GithubSourceReleaseAny),
}

var InstallGithubCommand = cli.Command{
	Name:    "github",
	Aliases: []string{"gh"},
	Usage:   "Install an application from Github",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "id",
			Usage: "Application identifier",
		},
		&cli.StringFlag{
			Name:    "release",
			Aliases: []string{"r"},
			Usage: fmt.Sprintf(
				"Releases type such as %s",
				strings.Join(githubSourceReleaseStrings, ", "),
			),
			Value: githubSourceReleaseStrings[0],
		},
		&cli.StringFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage: fmt.Sprintf(
				"Release tag name (requires release to be %s)",
				core.GithubSourceReleaseTagged,
			),
		},
		&cli.BoolFlag{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Creates a symlink",
		},
		&cli.BoolFlag{
			Name:    "assume-yes",
			Aliases: []string{"y"},
			Usage:   "Automatically answer yes for questions",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		utils.LogDebug("reading config")
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		args := cmd.Args()
		if args.Len() == 0 {
			return errors.New("no url specified")
		}
		if args.Len() > 1 {
			return errors.New("unexpected excessive arguments")
		}

		url := args.Get(0)
		appId := cmd.String("id")
		releaseType := cmd.String("release")
		tagName := cmd.String("tag")
		link := cmd.Bool("link")
		assumeYes := cmd.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument url: %s", url))
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))
		utils.LogDebug(fmt.Sprintf("argument release: %v", releaseType))
		utils.LogDebug(fmt.Sprintf("argument tag: %v", tagName))
		utils.LogDebug(fmt.Sprintf("argument link: %v", link))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		isValidUrl, ghUsername, ghReponame := core.ParseGithubRepoUrl(url)
		utils.LogDebug(fmt.Sprintf("parsed github url valid: %v", isValidUrl))
		utils.LogDebug(fmt.Sprintf("parsed github owner: %s", ghUsername))
		utils.LogDebug(fmt.Sprintf("parsed github repo: %s", ghReponame))
		if !isValidUrl {
			return errors.New("invalid github repo url")
		}
		if !utils.SliceContains(githubSourceReleaseStrings, releaseType) {
			return errors.New("invalid github release type")
		}

		if appId == "" {
			appId = core.ConstructAppId(ghReponame)
		}
		appId = utils.CleanId(appId)
		utils.LogDebug(fmt.Sprintf("clean id: %s", appId))
		if appId == "" {
			return errors.New("invalid application id")
		}

		source := &core.GithubSource{
			UserName: ghUsername,
			RepoName: ghReponame,
			Release:  core.GithubSourceRelease(releaseType),
			TagName:  tagName,
		}
		release, err := source.FetchAptRelease()
		if err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("selected github tag name: %s", release.TagName))

		matchScore, asset := release.ChooseAptAsset()
		if matchScore == core.AppImageAssetNoMatch {
			return fmt.Errorf("no valid asset in github tag %s", release.TagName)
		}
		if matchScore == core.AppImageAssetPartialMatch {
			utils.LogWarning("no architecture specified in the asset name, cannot determine compatibility")
		}
		utils.LogDebug(fmt.Sprintf("selected asset url %s", asset.DownloadUrl))

		appPaths := core.ConstructAppPaths(config, appId, &core.ConstructAppPathsOptions{
			Symlink: link,
		})
		if _, ok := config.Installed[appId]; ok {
			utils.LogWarning(fmt.Sprintf("application with id %s already exists", appId))
			if !assumeYes {
				proceed, err := utils.PromptYesNoInput(reader, "Do you want to re-install this application?")
				if err != nil {
					return err
				}
				if !proceed {
					utils.LogWarning("aborted...")
					return nil
				}
			}
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Identifier", color.CyanString(appId))
		summary.Add(utils.LogRightArrowPrefix, "Version", color.CyanString(release.TagName))
		summary.Add(utils.LogRightArrowPrefix, "Filename", color.CyanString(asset.Name))
		summary.Add(utils.LogRightArrowPrefix, "AppImage", color.CyanString(appPaths.AppImage))
		summary.Add(utils.LogRightArrowPrefix, ".desktop file", color.CyanString(appPaths.Desktop))
		if appPaths.Symlink != "" {
			summary.Add(utils.LogRightArrowPrefix, "Symlink", color.CyanString(appPaths.Symlink))
		}
		summary.Add(utils.LogRightArrowPrefix, "Download Size", color.CyanString(prettyBytes(asset.Size)))
		summary.Print()
		utils.LogLn()

		if !assumeYes {
			proceed, err := utils.PromptYesNoInput(reader, "Do you want to proceed?")
			if err != nil {
				return err
			}
			if !proceed {
				utils.LogWarning("aborted...")
				return nil
			}
		}

		app := &core.AppConfig{
			Id:      appId,
			Version: release.TagName,
			Source:  core.GithubSourceId,
			Paths:   *appPaths,
		}
		utils.LogLn()
		installed, _ := InstallApps([]InstallableApp{{
			App:    app,
			Source: source,
			Asset:  asset.ToAsset(),
		}})
		if installed != 1 {
			return nil
		}

		utils.LogLn()
		utils.LogInfo(
			fmt.Sprintf(
				"%s Installed %s successfully!",
				utils.LogTickPrefix,
				color.CyanString(app.Id),
			),
		)

		return nil
	},
}
