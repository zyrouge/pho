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
	},
	Action: func(ctx *cli.Context) error {
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		args := ctx.Args()

		appIds := args.Slice()
		assumeYes := ctx.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument ids: %s", strings.Join(appIds, ", ")))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		if len(appIds) == 0 {
			for x := range config.Installed {
				appIds = append(appIds, x)
			}
		}

		utils.LogLn()
		updateables, failed, err := CheckAppUpdates(config, appIds)
		if err != nil {
			return err
		}
		if failed > 0 {
			utils.LogLn()
		}
		if len(updateables) == 0 {
			utils.LogInfo(
				fmt.Sprintf(
					"%s Everything is up-to-date.",
					utils.LogTickPrefix,
				),
			)
			return nil
		}

		summary := utils.NewLogTable()
		headingColor := color.New(color.Underline, color.Bold)
		summary.Add(
			headingColor.Sprint("Index"),
			headingColor.Sprint("Application ID"),
			headingColor.Sprint("Application Name"),
			headingColor.Sprint("Old Version"),
			headingColor.Sprint("New Version"),
		)
		i := 0
		for _, x := range updateables {
			i++
			summary.Add(
				fmt.Sprintf("%d.", i),
				color.CyanString(x.App.Id),
				color.CyanString(x.App.Name),
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
				Paths:  x.Paths,
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
	Paths  *core.AppPaths
	Update *core.SourceUpdate
}

func CheckAppUpdates(config *core.Config, appIds []string) ([]UpdatableApp, int, error) {
	failed := 0
	apps := []UpdatableApp{}
	for _, appId := range appIds {
		updatable, err := CheckAppUpdate(config, appId)
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

func CheckAppUpdate(config *core.Config, appId string) (*UpdatableApp, error) {
	if _, ok := config.Installed[appId]; !ok {
		return nil, fmt.Errorf(
			"application with id %s is not installed",
			color.CyanString(appId),
		)
	}
	appPaths := core.ConstructAppPaths(config, appId, "")
	app, err := core.ReadAppConfig(appPaths.Config)
	if err != nil {
		return nil, err
	}
	appPaths = core.ConstructAppPaths(config, appId, app.Name)
	sourceConfig, err := core.ReadSourceConfig(app.Source, appPaths.SourceConfig)
	if err != nil {
		return nil, err
	}
	source, err := core.CastSourceConfigAsSource(sourceConfig)
	if err != nil {
		return nil, err
	}
	if !source.SupportUpdates() {
		return nil, nil
	}
	hasUpdate, update, err := source.CheckUpdate(app)
	if err != nil {
		return nil, err
	}
	if !hasUpdate {
		return nil, err
	}
	updatable := &UpdatableApp{
		App:    app,
		Source: sourceConfig,
		Paths:  appPaths,
		Update: update,
	}
	return updatable, nil
}
