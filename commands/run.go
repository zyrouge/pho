package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var RunCommand = cli.Command{
	Name:    "run",
	Aliases: []string{"launch"},
	Usage:   "Run an application",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "detached",
			Aliases: []string{"d"},
			Usage:   "Run as a detached process",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		utils.LogDebug("reading config")
		config, err := core.GetConfig()
		if err != nil {
			return err
		}

		args := cmd.Args()
		hasExecArgs := args.Get(1) == "--"
		if args.Len() == 0 {
			return errors.New("no application id specified")
		}
		if args.Len() > 1 && !hasExecArgs {
			return errors.New("unexpected excessive arguments")
		}

		appId := args.Get(0)
		execArgs := []string{}
		if hasExecArgs {
			execArgs = args.Slice()[2:]
		}
		detached := cmd.Bool("detached")
		utils.LogDebug(fmt.Sprintf("argument id: %s", appId))
		utils.LogDebug(fmt.Sprintf("argument exec-args: %s", strings.Join(execArgs, " ")))
		utils.LogDebug(fmt.Sprintf("argument detached: %v", detached))

		if _, ok := config.Installed[appId]; !ok {
			return fmt.Errorf(
				"application with id %s is not installed",
				color.CyanString(appId),
			)
		}

		appConfigPath := core.GetAppConfigPath(config, appId)
		utils.LogDebug(fmt.Sprintf("reading app config from %s", appConfigPath))
		app, err := core.ReadAppConfig(appConfigPath)
		if err != nil {
			return err
		}

		execPath := app.Paths.AppImage
		execDir, err := os.Getwd()
		if err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("exec path as %s", execPath))
		utils.LogDebug(fmt.Sprintf("exec dir as %s", execDir))

		if detached {
			detachedOptions := &utils.StartDetachedProcessOptions{
				Dir:  execDir,
				Exec: execPath,
				Args: execArgs,
			}
			if err = utils.StartDetachedProcess(detachedOptions); err != nil {
				return err
			}
			utils.LogDebug("launched detached process successfully")
			return nil
		}

		proc := exec.Command(execPath)
		proc.Dir = execDir
		proc.Stdin = os.Stdin
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		if err = proc.Run(); err != nil {
			return err
		}
		return nil
	},
}
