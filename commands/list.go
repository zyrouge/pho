package commands

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var ListCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"installed"},
	Usage:   "List all installed applications",
	Action: func(_ context.Context, cmd *cli.Command) error {
		utils.LogDebug("reading config")
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		headingColor := color.New(color.Underline, color.Bold)
		summary.Add(
			headingColor.Sprint("Index"),
			headingColor.Sprint("Application ID"),
		)
		i := 0
		for appId := range config.Installed {
			i++
			summary.Add(fmt.Sprintf("%d.", i), color.CyanString(appId))
		}
		summary.Print()
		if i == 0 {
			utils.LogInfo(color.HiBlackString("no applications are installed"))
		}
		utils.LogLn()

		return nil
	},
}
