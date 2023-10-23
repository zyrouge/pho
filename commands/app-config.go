package commands

import (
	"github.com/urfave/cli/v3"
)

var AppConfigCommand = cli.Command{
	Name:    "app-config",
	Aliases: []string{"ac"},
	Usage:   "Related to application configuration",
	Commands: []*cli.Command{
		&AppConfigSetIdCommand,
	},
}
