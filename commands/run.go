package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var RunCommand = cli.Command{
	Name:    "run",
	Aliases: []string{"open", "launch"},
	Usage:   "Run an application",
	Action: func(ctx *cli.Context) error {
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		args := ctx.Args()
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

		appConfigPath := core.ConstructAppConfigPath(config, appId)
		app, err := core.ReadAppConfig(appConfigPath)
		if err != nil {
			return err
		}

		stdin, err := os.Open(os.DevNull)
		if err != nil {
			return err
		}
		stdout, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		stderr, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return err
		}

		procAttr := &os.ProcAttr{
			Dir: path.Dir(app.Paths.AppImage),
			Env: os.Environ(),
			Files: []*os.File{
				stdin,
				stdout,
				stderr,
			},
			Sys: &syscall.SysProcAttr{
				Foreground: true,
			},
		}
		proc, err := os.StartProcess(app.Paths.AppImage, []string{}, procAttr)
		if err != nil {
			return err
		}
		if err = proc.Release(); err != nil {
			return err
		}

		utils.LogLn()
		utils.LogInfo(
			fmt.Sprintf(
				"%s Launched %s successfully!",
				utils.LogTickPrefix,
				color.CyanString(app.Id),
			),
		)

		return nil
	},
}
