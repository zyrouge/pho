package commands

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var AppConfigSetIdCommand = cli.Command{
	Name:  "set-id",
	Usage: "Update an application's identifier",
	Action: func(ctx *cli.Context) error {
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		args := ctx.Args()
		if args.Len() == 0 {
			return errors.New("no application id specified")
		}
		if args.Len() == 1 {
			return errors.New("no new application id specified")
		}
		if args.Len() > 2 {
			return errors.New("unexpected excessive arguments")
		}

		fromAppId := args.Get(0)
		if _, ok := config.Installed[fromAppId]; !ok {
			return fmt.Errorf(
				"application with id %s is not installed",
				color.CyanString(fromAppId),
			)
		}

		toAppId := core.ConstructAppId(args.Get(1))
		if toAppId == "" {
			return errors.New("invalid application id")
		}

		fromAppConfigPath := core.GetAppConfigPath(config, fromAppId)
		app, err := core.ReadAppConfig(fromAppConfigPath)
		if err != nil {
			return err
		}
		fromAppPaths := app.Paths
		toAppPaths := core.ConstructAppPaths(config, toAppId)
		app.Id = toAppId
		app.Paths = *toAppPaths
		delete(config.Installed, fromAppId)
		config.Installed[toAppId] = toAppPaths.Config
		if err = os.Rename(fromAppPaths.Dir, toAppPaths.Dir); err != nil {
			return err
		}
		fromAppImagePath := path.Join(toAppPaths.Dir, path.Base(fromAppPaths.AppImage))
		if err = os.Rename(fromAppImagePath, toAppPaths.AppImage); err != nil {
			return err
		}
		fromIconPath := path.Join(toAppPaths.Dir, path.Base(fromAppPaths.Icon))
		if err = os.Rename(fromIconPath, toAppPaths.Icon); err != nil {
			return err
		}
		if err = core.SaveAppConfig(toAppPaths.Config, app); err != nil {
			return err
		}
		if err = core.SaveConfig(config); err != nil {
			return err
		}
		desktopContent, err := os.ReadFile(fromAppPaths.Desktop)
		if err != nil {
			return err
		}
		if err = core.UninstallDesktopFile(fromAppPaths.Desktop); err != nil {
			return err
		}
		if err = core.InstallDesktopFile(toAppPaths, string(desktopContent)); err != nil {
			return err
		}

		utils.LogLn()
		utils.LogInfo(
			fmt.Sprintf(
				"%s Renamed %s to %s successfully!",
				utils.LogTickPrefix,
				color.CyanString(fromAppId),
				color.CyanString(toAppId),
			),
		)

		return nil
	},
}
