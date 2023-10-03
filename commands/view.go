package commands

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var ViewCommand = cli.Command{
	Name:    "view",
	Aliases: []string{},
	Usage:   "View an installed application",
	Action: func(ctx *cli.Context) error {
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		args := ctx.Args()
		if args.Len() == 0 {
			return errors.New("no application id specified")
		}
		if args.Len() > 1 {
			return errors.New("unexpected excessive arguments")
		}

		appId := args.Get(0)
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))

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
		appPaths = core.GetAppPaths(config, appId, app.Name)

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Name", color.CyanString(app.Name))
		summary.Add(utils.LogRightArrowPrefix, "Identifier", color.CyanString(app.Id))
		summary.Add(utils.LogRightArrowPrefix, "Version", color.CyanString(app.Version))
		summary.Add(utils.LogRightArrowPrefix, "Directory", color.CyanString(appPaths.Dir))
		summary.Add(utils.LogRightArrowPrefix, "AppImage", color.CyanString(app.AppImage))
		summary.Add(utils.LogRightArrowPrefix, "Icon", color.CyanString(appPaths.Icon))
		summary.Add(utils.LogRightArrowPrefix, ".desktop file", color.CyanString(appPaths.Desktop))
		summary.Add(utils.LogRightArrowPrefix, "Source", color.CyanString(string(app.Source)))
		summary.Print()
		utils.LogLn()

		return nil
	},
}