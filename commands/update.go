package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		&cli.BoolFlag{
			Name:  "reinstall",
			Usage: "Forcefully update and reinstall application",
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

		appIds := args.Slice()
		assumeYes := ctx.Bool("assume-yes")
		reinstall := ctx.Bool("reinstall")
		utils.LogDebug(fmt.Sprintf("argument ids: %s", strings.Join(appIds, ", ")))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))
		utils.LogDebug(fmt.Sprintf("argument reinstall: %v", reinstall))

		if len(appIds) == 0 {
			for x := range config.Installed {
				appIds = append(appIds, x)
			}
		}

		updateables, _, err := CheckAppUpdates(config, appIds, reinstall)
		if err != nil {
			return err
		}
		if len(updateables) == 0 {
			utils.LogLn()
			utils.LogInfo(
				fmt.Sprintf(
					"%s Everything is up-to-date.",
					utils.LogTickPrefix,
				),
			)
			return nil
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		headingColor := color.New(color.Underline, color.Bold)
		summary.Add(
			headingColor.Sprint("Index"),
			headingColor.Sprint("Application ID"),
			headingColor.Sprint("Old Version"),
			headingColor.Sprint("New Version"),
		)
		i := 0
		for _, x := range updateables {
			i++
			summary.Add(
				fmt.Sprintf("%d.", i),
				color.CyanString(x.App.Id),
				x.App.Version,
				color.CyanString(x.Update.Version),
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
		installables := []InstallableApp{}
		for _, x := range updateables {
			x.App.Version = x.Update.Version
			installables = append(installables, InstallableApp{
				App:    x.App,
				Source: x.Source,
				Asset:  x.Update.Asset,
			})
		}
		installed, failed := InstallApps(installables)

		utils.LogLn()
		if installed > 0 {
			utils.LogInfo(
				fmt.Sprintf(
					"%s Updated %s applications successfully!",
					utils.LogTickPrefix,
					color.CyanString(fmt.Sprint(installed)),
				),
			)
		}
		if failed > 0 {
			utils.LogInfo(
				fmt.Sprintf(
					"%s Failed to update %s applications.",
					utils.LogExclamationPrefix,
					color.RedString(fmt.Sprint(failed)),
				),
			)
		}

		return nil
	},
}

type UpdatableApp struct {
	App    *core.AppConfig
	Source any
	Update *core.SourceUpdate
}

func CheckAppUpdates(config *core.Config, appIds []string, reinstall bool) ([]UpdatableApp, int, error) {
	failed := 0
	apps := []UpdatableApp{}
	for _, appId := range appIds {
		updatable, err := CheckAppUpdate(config, appId, reinstall)
		if err != nil {
			failed++
			utils.LogError(err)
			continue
		}
		if updatable != nil {
			apps = append(apps, *updatable)
		}
	}
	return apps, failed, nil
}

func CheckAppUpdate(config *core.Config, appId string, reinstall bool) (*UpdatableApp, error) {
	if _, ok := config.Installed[appId]; !ok {
		return nil, fmt.Errorf(
			"application with id %s is not installed",
			color.CyanString(appId),
		)
	}
	appConfigPath := core.GetAppConfigPath(config, appId)
	utils.LogDebug(fmt.Sprintf("reading app config from %s", appConfigPath))
	app, err := core.ReadAppConfig(appConfigPath)
	if err != nil {
		return nil, err
	}
	utils.LogDebug(fmt.Sprintf("reading app source config from %s", app.Paths.SourceConfig))
	sourceConfig, err := core.ReadSourceConfig(app.Source, app.Paths.SourceConfig)
	if err != nil {
		return nil, err
	}
	source, err := core.CastSourceConfigAsSource(sourceConfig)
	if err != nil {
		return nil, err
	}
	if !source.SupportUpdates() {
		utils.LogDebug(fmt.Sprintf("%s doesnt support any updates", appId))
		return nil, nil
	}
	update, err := source.CheckUpdate(app, reinstall)
	if err != nil {
		return nil, err
	}
	if update == nil {
		utils.LogDebug(fmt.Sprintf("%s has no updates", appId))
		return nil, nil
	}
	updatable := &UpdatableApp{
		App:    app,
		Source: sourceConfig,
		Update: update,
	}
	return updatable, nil
}
