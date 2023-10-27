package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var InstallLocalCommand = cli.Command{
	Name:  "local",
	Usage: "Install local AppImage",
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
	Action: func(ctx *cli.Context) error {
		utils.LogDebug("reading config")
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

		appImagePath := args.Get(0)
		appId := ctx.String("id")
		appVersion := ctx.String("version")
		assumeYes := ctx.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument path: %s", appImagePath))
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		if appImagePath == "" {
			return errors.New("invalid appimage path")
		}
		if !path.IsAbs(appImagePath) {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			appImagePath = path.Join(cwd, appImagePath)
		}
		utils.LogDebug(fmt.Sprintf("resolved appimage path: %s", appImagePath))
		appImageFileInfo, err := os.Stat(appImagePath)
		if err != nil {
			return err
		}
		if appImageFileInfo.IsDir() {
			return errors.New("appimage path must be a file")
		}

		if appId == "" {
			appId = core.ConstructAppId(path.Base(appImagePath))
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

		app := &core.AppConfig{
			Id:      appId,
			Version: appVersion,
			Source:  core.LocalSourceId,
			Paths:   *appPaths,
		}
		source := &core.LocalSource{}
		asset := &core.Asset{
			Source:   appImagePath,
			Size:     appImageFileInfo.Size(),
			Download: core.LocalAssetDownload(appImagePath),
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
