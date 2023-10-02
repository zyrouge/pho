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

var UninstallCommand = cli.Command{
	Name:    "uninstall",
	Aliases: []string{"remove"},
	Usage:   "Uninstall an application",
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

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Name", color.RedString(app.Name))
		summary.Add(utils.LogRightArrowPrefix, "Identifier", color.RedString(app.Id))
		summary.Add(utils.LogRightArrowPrefix, "Version", color.RedString(app.Version))
		summary.Add(utils.LogRightArrowPrefix, "AppImage", color.RedString(app.AppImage))
		summary.Add(utils.LogRightArrowPrefix, ".desktop file", color.RedString(appPaths.Desktop))
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
		errorCount := UninstallApp(app, appPaths)
		if errorCount > 0 {
			utils.LogLn()
			utils.LogInfo(
				fmt.Sprintf(
					"%s Uninstalled %s with %d errors.",
					utils.LogTickPrefix,
					color.RedString(app.Name),
					errorCount,
				),
			)
		} else {
			utils.LogInfo(
				fmt.Sprintf(
					"%s Uninstalled %s successfully!",
					utils.LogTickPrefix,
					color.RedString(app.Name),
				),
			)
		}

		return nil
	},
}

func UninstallApp(app *core.AppConfig, appPaths *core.AppPaths) int {
	failed := 0
	config, err := core.ReadConfig()
	if err != nil {
		utils.LogError(err)
		failed++
	} else {
		delete(config.Installed, app.Id)
		if err = core.SaveConfig(config); err != nil {
			utils.LogError(err)
			failed++
		}
	}
	if err = os.RemoveAll(appPaths.Dir); err != nil {
		utils.LogError(err)
		failed++
	}
	if err = core.UninstallDesktopFile(appPaths.Desktop); err != nil {
		utils.LogError(err)
		failed++
	}
	if err = os.Remove(appPaths.Desktop); err != nil {
		utils.LogError(err)
		failed++
	}
	return failed
}
