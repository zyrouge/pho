package commands

import (
	"bufio"
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

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "Initialize and setup necessities",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "apps-dir",
			Usage: "AppImages directory",
		},
		&cli.StringFlag{
			Name:  "apps-desktop-dir",
			Usage: ".desktop files directory",
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite config if exists",
		},
		&cli.BoolFlag{
			Name:    "assume-yes",
			Aliases: []string{"y"},
			Usage:   "Automatically answer yes for questions",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		appsDir := cmd.String("apps-dir")
		appsDesktopDir := cmd.String("apps-desktop-dir")
		overwrite := cmd.Bool("overwrite")
		assumeYes := cmd.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument apps-dir: %s", appsDir))
		utils.LogDebug(fmt.Sprintf("argument apps-desktop-dir: %s", appsDesktopDir))
		utils.LogDebug(fmt.Sprintf("argument overwrite: %v", overwrite))
		utils.LogDebug(fmt.Sprintf("argument assume-yes: %v", assumeYes))

		reader := bufio.NewReader(os.Stdin)
		configPath, err := core.GetConfigPath()
		if err != nil {
			return err
		}
		configExists, err := utils.FileExists(configPath)
		if err != nil {
			return err
		}
		if configExists {
			utils.LogWarning("config already exists")
			if !overwrite {
				if assumeYes {
					return fmt.Errorf(
						"pass in %s flag to overwrite configuration file",
						color.CyanString("--overwrite"),
					)
				}
				proceed, err := utils.PromptYesNoInput(
					reader,
					"Do you want to re-initiliaze configuration file?",
				)
				if err != nil {
					return err
				}
				if !proceed {
					return nil
				}
			}
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		if appsDir == "" {
			appsDir = path.Join(homeDir, ".local/share", core.AppCodeName, "applications")
			if !assumeYes {
				appsDir, err = utils.PromptTextInput(
					reader,
					"Where do you want to store the AppImages?",
					appsDir,
				)
				if err != nil {
					return err
				}
			}
		}
		if appsDir == "" {
			return errors.New("invalid application name")
		}
		appsDir, err = utils.ResolvePath(appsDir)
		if err != nil {
			return err
		}

		if appsDesktopDir == "" {
			appsDesktopDir = path.Join(homeDir, ".local/share", "applications")
			if !assumeYes {
				appsDesktopDir, err = utils.PromptTextInput(
					reader,
					"Where do you want to store the .desktop files?",
					appsDesktopDir,
				)
				if err != nil {
					return err
				}
			}
		}
		if appsDesktopDir == "" {
			return errors.New("invalid application desktop file path")
		}
		appsDesktopDir, err = utils.ResolvePath(appsDesktopDir)
		if err != nil {
			return err
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Configuration file", color.CyanString(configPath))
		summary.Add(utils.LogRightArrowPrefix, "AppImages directory", color.CyanString(appsDir))
		summary.Add(utils.LogRightArrowPrefix, ".desktop files directory", color.CyanString(appsDesktopDir))
		summary.Print()
		utils.LogLn()

		if !assumeYes {
			proceed, err := utils.PromptYesNoInput(reader, "Do you want to proceed?")
			if err != nil {
				return err
			}
			if !proceed {
				utils.LogWarning("aborted...")
				return nil
			}
		}

		if err := os.MkdirAll(path.Dir(configPath), os.ModePerm); err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(appsDir), os.ModePerm); err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(appsDesktopDir), os.ModePerm); err != nil {
			return err
		}
		config := &core.Config{
			AppsDir:    appsDir,
			DesktopDir: appsDesktopDir,
			Installed:  map[string]string{},
		}
		err = core.SaveConfig(config)
		if err != nil {
			return err
		}
		utils.LogInfo(fmt.Sprintf("%s Generated %s", utils.LogTickPrefix, color.CyanString(configPath)))

		return nil
	},
}
