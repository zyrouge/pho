package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

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
			Name:  "tag",
			Usage: "Tag name",
		},
		&cli.BoolFlag{
			Name:  "prerelease",
			Usage: "Select pre-release tags",
			Value: false,
		},
		&cli.BoolFlag{
			Name:    "assume-yes",
			Aliases: []string{"y"},
			Usage:   "Automatically answer yes for questions",
		},
	},
	Action: func(ctx *cli.Context) error {
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		args := ctx.Args()
		if args.Len() == 0 {
			return errors.New("no url specified")
		}
		if args.Len() > 1 {
			return errors.New("unexpected excessive arguments")
		}

		url := args.Get(0)
		appId := ctx.String("id")
		tagName := ctx.String("tag")
		prerelease := ctx.Bool("prerelease")
		assumeYes := ctx.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument url: %s", url))
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))
		utils.LogDebug(fmt.Sprintf("argument tag: %s", tagName))
		utils.LogDebug(fmt.Sprintf("argument prerelease: %v", prerelease))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		isValidUrl, ghUsername, ghReponame := core.ParseGithubRepoUrl(url)
		utils.LogDebug(fmt.Sprintf("parsed github url valid: %v", isValidUrl))
		utils.LogDebug(fmt.Sprintf("parsed github owner: %s", ghUsername))
		utils.LogDebug(fmt.Sprintf("parsed github repo: %s", ghReponame))
		if !isValidUrl {
			return errors.New("invalid github repo url")
		}

		if appId == "" {
			appId = core.ConstructAppId(fmt.Sprintf("%s-%s", ghUsername, ghReponame))
		}
		appId = utils.CleanId(appId)
		if appId == "" {
			return errors.New("invalid application id")
		}

		source := &core.GithubSource{
			UserName:   ghUsername,
			RepoName:   ghReponame,
			PreRelease: prerelease,
			TagName:    tagName,
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

		appPaths := core.ConstructAppPaths(config, appId)
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
