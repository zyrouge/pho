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

var UpdateCommand = cli.Command{
	Name:    "update",
	Aliases: []string{"upgrade"},
	Usage:   "Update an application",
	Flags: []cli.Flag{
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
			return errors.New("no application id specified")
		}
		if args.Len() > 1 {
			return errors.New("unexpected excessive arguments")
		}

		appId := args.Get(0)
		assumeYes := ctx.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		if _, ok := config.Installed[appId]; !ok {
			return fmt.Errorf(
				"application with id %s is not installed",
				color.CyanString(appId),
			)
		}

		appPaths := core.GetAppPaths(config, appId, "")
		app, err := core.ReadAppConfig(appPaths.Config)
		if err != nil {
			return err
		}
		// re-constrct paths since the previous one did not specify app name
		appPaths = core.GetAppPaths(config, appId, app.Name)
		sourceConfig, err := core.ReadSourceConfig(app.Source, appPaths.SourceConfig)
		if err != nil {
			return err
		}
		source, err := core.CastSourceConfigAsSource(sourceConfig)
		if err != nil {
			return err
		}
		if !source.SupportUpdates() {
			return errors.New("application does not support automatic updates")
		}
		hasUpdates, update, err := source.CheckUpdate(app)
		if err != nil {
			return err
		}
		if !hasUpdates {
			utils.LogLn()
			utils.LogInfo(
				fmt.Sprintf(
					"%s %s is already up-to-date.",
					utils.LogTickPrefix,
					color.CyanString(app.Name),
				),
			)
			return nil
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Name", color.CyanString(app.Name))
		summary.Add(utils.LogRightArrowPrefix, "Identifier", color.CyanString(app.Id))
		summary.Add(
			utils.LogRightArrowPrefix,
			"Version",
			fmt.Sprintf(
				"%s %s %s",
				color.HiBlackString(app.Version),
				color.HiBlackString("->"),
				color.CyanString(update.Version),
			),
		)
		summary.Add(utils.LogRightArrowPrefix, "AppImage", color.CyanString(app.AppImage))
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

		utils.LogLn()
		installed, _ := InstallApps([]InstallableApp{{
			App:    app,
			Source: sourceConfig,
			Paths:  appPaths,
			Asset:  update.Asset,
		}})
		if installed != 1 {
			return nil
		}

		utils.LogLn()
		utils.LogInfo(
			fmt.Sprintf(
				"%s Updated %s successfully!",
				utils.LogTickPrefix,
				color.CyanString(app.Name),
			),
		)

		return nil
	},
}
