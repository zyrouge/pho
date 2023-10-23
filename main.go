package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/commands"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

func main() {
	app := &cli.Command{
		Name:    core.AppExecutableName,
		Usage:   core.AppDescription,
		Version: core.AppVersion,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:       "verbose",
				Persistent: true,
				Action: func(ctx *cli.Context, b bool) error {
					utils.LogDebugEnabled = b
					return nil
				},
			},
		},
		Authors: []any{"Zyrouge"},
		Commands: []*cli.Command{
			&commands.InitCommand,
			&commands.InstallCommand,
			&commands.UninstallCommand,
			&commands.UpdateCommand,
			&commands.RunCommand,
			&commands.ListCommand,
			&commands.ViewCommand,
			&commands.TidyBrokenCommand,
			&commands.SelfUpdateCommand,
			&commands.AppConfigCommand,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		utils.LogError(err)
		os.Exit(1)
	}
}
