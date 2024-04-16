package commands

import (
	"context"
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
	Action: func(_ context.Context, cmd *cli.Command) error {
		utils.LogDebug("reading config")
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		args := cmd.Args()
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

		appConfigPath := core.GetAppConfigPath(config, appId)
		utils.LogDebug(fmt.Sprintf("reading app config from %s", appConfigPath))
		app, err := core.ReadAppConfig(appConfigPath)
		if err != nil {
			return err
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Identifier", color.CyanString(app.Id))
		summary.Add(utils.LogRightArrowPrefix, "Version", color.CyanString(app.Version))
		summary.Add(utils.LogRightArrowPrefix, "Source", color.CyanString(string(app.Source)))
		summary.Add(utils.LogRightArrowPrefix, "Directory", color.CyanString(app.Paths.Dir))
		summary.Add(utils.LogRightArrowPrefix, "AppImage", color.CyanString(app.Paths.AppImage))
		summary.Add(utils.LogRightArrowPrefix, "Icon", color.CyanString(app.Paths.Icon))
		summary.Add(utils.LogRightArrowPrefix, ".desktop file", color.CyanString(app.Paths.Desktop))
		summary.Print()
		utils.LogLn()

		return nil
	},
}
