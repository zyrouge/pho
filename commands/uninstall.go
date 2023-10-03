package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var UninstallCommand = cli.Command{
	Name:    "uninstall",
	Aliases: []string{"remove", "delete"},
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
			return errors.New("no application ids specified")
		}

		appIds := args.Slice()
		assumeYes := ctx.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument ids: %s", strings.Join(appIds, ", ")))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		utils.LogLn()
		failed := 0
		uninstallables := []core.AppConfig{}
		for _, appId := range appIds {
			if _, ok := config.Installed[appId]; !ok {
				failed++
				utils.LogError(
					fmt.Sprintf(
						"application with id %s is not installed",
						color.CyanString(appId),
					),
				)
				continue
			}
			appPaths := core.GetAppPaths(config, appId, "")
			app, err := core.ReadAppConfig(appPaths.Config)
			if err != nil {
				failed++
				utils.LogError(err)
				continue
			}
			uninstallables = append(uninstallables, *app)
		}
		if len(uninstallables) == 0 {
			return nil
		}
		if failed > 0 {
			utils.LogLn()
		}

		summary := utils.NewLogTable()
		headingColor := color.New(color.Underline, color.Bold)
		summary.Add(
			headingColor.Sprint("Index"),
			headingColor.Sprint("Application ID"),
			headingColor.Sprint("Application Name"),
			headingColor.Sprint("Version"),
		)
		i := 0
		for _, x := range uninstallables {
			i++
			summary.Add(
				fmt.Sprintf("%d.", i),
				color.RedString(x.Id),
				color.RedString(x.Name),
				x.Version,
			)
		}
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
		failed = 0
		for _, x := range uninstallables {
			appPaths := core.GetAppPaths(config, x.Id, x.Name)
			failed += UninstallApp(&x, appPaths)
		}
		if failed > 0 {
			utils.LogLn()
			utils.LogInfo(
				fmt.Sprintf(
					"%s Uninstalled %s applications with %s errors.",
					utils.LogTickPrefix,
					color.RedString(fmt.Sprint(len(uninstallables))),
					color.RedString(fmt.Sprint(failed)),
				),
			)
		} else {
			utils.LogInfo(
				fmt.Sprintf(
					"%s Uninstalled %s applications successfully!",
					utils.LogTickPrefix,
					color.RedString(fmt.Sprint(len(uninstallables))),
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
	utils.LogDebug(fmt.Sprintf("removing %s", appPaths.Dir))
	if err = os.RemoveAll(appPaths.Dir); err != nil {
		utils.LogError(err)
		failed++
	}
	utils.LogDebug(fmt.Sprintf("removing %s", appPaths.Desktop))
	if err = core.UninstallDesktopFile(appPaths.Desktop); err != nil {
		utils.LogError(err)
		failed++
	}
	utils.LogDebug(fmt.Sprintf("removing %s", appPaths.Desktop))
	if err = os.Remove(appPaths.Desktop); err != nil {
		utils.LogError(err)
		failed++
	}
	return failed
}
