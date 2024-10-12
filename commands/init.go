package commands

import (
	"bufio"
	"context"
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
		&cli.StringFlag{
			Name:  "apps-link-dir",
			Usage: "AppImage symlinks directory",
		},
		&cli.BoolFlag{
			Name:  "enable-integration-prompt",
			Usage: "Enables AppImageLauncher's integration prompt",
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite config if exists",
		},
		&cli.BoolFlag{
			Name:    "assume-yes",
			Aliases: []string{"y"},
			Usage:   "Automatically answer 'yes' for questions",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		appsDir := cmd.String("apps-dir")
		appsDesktopDir := cmd.String("apps-desktop-dir")
		appsLinkDir := cmd.String("apps-link-dir")
		enableIntegrationPromptSet, enableIntegrationPrompt := utils.CommandBoolSetAndValue(cmd, "enable-integration-prompt")
		overwrite := cmd.Bool("overwrite")
		assumeYes := cmd.Bool("assume-yes")
		utils.LogDebug(fmt.Sprintf("argument apps-dir: %s", appsDir))
		utils.LogDebug(fmt.Sprintf("argument apps-desktop-dir: %s", appsDesktopDir))
		utils.LogDebug(fmt.Sprintf("argument apps-link-dir: %s", appsLinkDir))
		utils.LogDebug(fmt.Sprintf("argument enable-integration-prompt: %v", enableIntegrationPrompt))
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
		appsDesktopDir, err = utils.ResolvePath(appsDesktopDir)
		if err != nil {
			return err
		}

		enableAppsLinkDir := true
		if !assumeYes {
			enableAppsLinkDir, err = utils.PromptYesNoInput(
				reader,
				"Do you want to symlink AppImage files?",
			)
			if err != nil {
				return err
			}
		}
		if enableAppsLinkDir && appsLinkDir == "" {
			appsLinkDir = path.Join(homeDir, ".local/bin")
			if !assumeYes {
				appsLinkDir, err = utils.PromptTextInput(
					reader,
					"Where do you want to symlink AppImage files?",
					appsLinkDir,
				)
				if err != nil {
					return err
				}
			}
		}
		if appsLinkDir != "" {
			appsLinkDir, err = utils.ResolvePath(appsLinkDir)
			if err != nil {
				return err
			}
		}

		if !enableIntegrationPromptSet && !assumeYes {
			enableIntegrationPrompt, err = utils.PromptYesNoInput(
				reader,
				"Do you want to enable AppImageLauncher's integration prompt?",
			)
			if err != nil {
				return err
			}
		}

		utils.LogLn()
		summary := utils.NewLogTable()
		summary.Add(utils.LogRightArrowPrefix, "Configuration file", color.CyanString(configPath))
		summary.Add(utils.LogRightArrowPrefix, "AppImages directory", color.CyanString(appsDir))
		summary.Add(utils.LogRightArrowPrefix, ".desktop files directory", color.CyanString(appsDesktopDir))
		if enableAppsLinkDir {
			summary.Add(utils.LogRightArrowPrefix, "AppImages symlink directory", color.CyanString(appsLinkDir))
		}
		summary.Add(utils.LogRightArrowPrefix, "Enable AppImageLauncher's integration prompt?", color.CyanString(utils.BoolToYesNo(enableIntegrationPrompt)))
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
			AppsDir:                 appsDir,
			DesktopDir:              appsDesktopDir,
			Installed:               map[string]string{},
			EnableIntegrationPrompt: enableIntegrationPrompt,
			SymlinksDir:             appsLinkDir,
		}
		err = core.SaveConfig(config)
		if err != nil {
			return err
		}
		utils.LogInfo(fmt.Sprintf("%s Generated %s", utils.LogTickPrefix, color.CyanString(configPath)))

		return nil
	},
}
