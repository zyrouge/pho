package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var ListCommand = cli.Command{
	Name:    "list",
	Aliases: []string{},
	Usage:   "List all installed applications",
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

		utils.LogLn()
		summary := utils.NewLogTable()
		headingColor := color.New(color.Underline, color.Bold)
		summary.Add(
			headingColor.Sprint("Index"),
			headingColor.Sprint("Application ID"),
			headingColor.Sprint("Installed Path"),
		)
		i := 0
		for appId, appDir := range config.Installed {
			i++
			summary.Add(fmt.Sprintf("%d.", i), color.CyanString(appId), appDir)
		}
		summary.Print()
		if i == 0 {
			utils.LogInfo(color.HiBlackString("no applications are installed"))
		}
		utils.LogLn()

		return nil
	},
}
