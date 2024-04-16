package commands

import (
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

var AppConfigSetIdCommand = cli.Command{
	Name:  "set-id",
	Usage: "Update an application's identifier",
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
		if args.Len() == 1 {
			return errors.New("no new application id specified")
		}
		if args.Len() > 2 {
			return errors.New("unexpected excessive arguments")
		}

		fromAppId := args.Get(0)
		toAppId := args.Get(1)
		utils.LogDebug(fmt.Sprintf("argument from-id: %v", fromAppId))
		utils.LogDebug(fmt.Sprintf("argument to-id: %v", toAppId))

		if _, ok := config.Installed[fromAppId]; !ok {
			return fmt.Errorf(
				"application with id %s is not installed",
				color.CyanString(fromAppId),
			)
		}

		toAppId = core.ConstructAppId(toAppId)
		utils.LogDebug(fmt.Sprintf("clean to-id: %v", toAppId))
		if toAppId == "" {
			return errors.New("invalid application id")
		}

		fromAppConfigPath := core.GetAppConfigPath(config, fromAppId)
		utils.LogDebug(fmt.Sprintf("reading app config from %s", fromAppConfigPath))
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
		utils.LogDebug(fmt.Sprintf("moving from %s to %s", fromAppPaths.Dir, toAppPaths.Dir))
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
		utils.LogDebug(fmt.Sprintf("saving app config to %s", toAppPaths.Config))
		if err = core.SaveAppConfig(toAppPaths.Config, app); err != nil {
			return err
		}
		utils.LogDebug("saving config")
		if err = core.SaveConfig(config); err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("reading .desktop file at %s", fromAppPaths.Desktop))
		desktopContent, err := os.ReadFile(fromAppPaths.Desktop)
		if err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("uninstalling .desktop file at %s", fromAppPaths.Desktop))
		if err = core.UninstallDesktopFile(fromAppPaths.Desktop); err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("installing .desktop file at %s", fromAppPaths.Desktop))
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
