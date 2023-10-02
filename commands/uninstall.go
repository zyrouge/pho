package commands

import "github.com/urfave/cli/v3"

var UninstallCommand = cli.Command{
	Name:    "uninstall",
	Aliases: []string{"remove"},
	Usage:   "Uninstall an application",
	Action: func(ctx *cli.Context) error {
		return nil
	},
}
