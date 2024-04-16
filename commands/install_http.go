package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var InstallHttpCommand = cli.Command{
	Name:    "http",
	Aliases: []string{"network", "from-url"},
	Usage:   "Install AppImage from http url",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "id",
			Usage: "Application identifier",
		},
		&cli.StringFlag{
			Name:  "version",
			Usage: "Application version",
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
		appVersion := cmd.String("version")
		assumeYes := cmd.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument url: %s", url))
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		if url == "" {
			return errors.New("invalid url")
		}

		if appId == "" {
			appId = core.ConstructAppId(path.Base(url))
			if !assumeYes {
				appId, err = utils.PromptTextInput(
					reader,
					"What should be the Application ID?",
					appId,
				)
				if err != nil {
					return err
				}
			}
		}
		appId = utils.CleanId(appId)
		utils.LogDebug(fmt.Sprintf("clean id: %s", appId))
		if appId == "" {
			return errors.New("invalid application id")
		}

		if appVersion == "" {
			appVersion = "0.0.0"
		}

		appPaths := core.ConstructAppPaths(config, appId)
		if _, ok := config.Installed[appId]; ok {
			utils.LogWarning(
				fmt.Sprintf(
					"application with id %s already exists",
					color.CyanString(appId),
				),
			)
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
		summary.Add(utils.LogRightArrowPrefix, "Version", color.CyanString(appVersion))
		summary.Add(utils.LogRightArrowPrefix, "AppImage", color.CyanString(appPaths.AppImage))
		summary.Add(utils.LogRightArrowPrefix, ".desktop file", color.CyanString(appPaths.Desktop))
		summary.Print()

		if !assumeYes {
			utils.LogLn()
			proceed, err := utils.PromptYesNoInput(reader, "Do you want to proceed?")
			if err != nil {
				return err
			}
			if !proceed {
				utils.LogWarning("aborted...")
				return nil
			}
		}

		assetMetadata, err := core.ExtractNetworkAssetMetadata(url)
		if err != nil {
			return err
		}

		app := &core.AppConfig{
			Id:      appId,
			Version: appVersion,
			Source:  core.HttpSourceId,
			Paths:   *appPaths,
		}
		source := &core.HttpSource{}
		asset := &core.Asset{
			Source:   url,
			Size:     assetMetadata.Size,
			Download: core.NetworkAssetDownload(url),
		}

		utils.LogLn()
		installed, _ := InstallApps([]InstallableApp{{
			App:    app,
			Source: source,
			Asset:  asset,
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
