package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var SelfUpdateCommand = cli.Command{
	Name:    "self-update",
	Aliases: []string{"self-upgrade"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "reinstall",
			Usage: "Forcefully update",
		},
	},
	Usage: fmt.Sprintf("Update %s", core.AppName),
	Action: func(ctx *cli.Context) error {
		reinstall := ctx.Bool("reinstall")
		utils.LogDebug(fmt.Sprintf("argument reinstall: %v", reinstall))

		utils.LogDebug("fetching latest release")
		release, err := core.GithubApiFetchLatestRelease(core.AppGithubOwner, core.AppGithubRepo)
		if err != nil {
			return err
		}
		if release.TagName == fmt.Sprintf("v%s", core.AppVersion) && !reinstall {
			utils.LogInfo(
				fmt.Sprintf("%s You are already on the latest version!", utils.LogTickPrefix),
			)
			return nil
		}
		arch := utils.GetSystemArch()
		var asset *core.GithubApiReleaseAsset
		for i := range release.Assets {
			x := release.Assets[i]
			if strings.HasSuffix(x.Name, arch) {
				asset = &x
				break
			}
		}
		if asset == nil {
			return fmt.Errorf(
				"unable to find appropriate binary from release %s",
				release.TagName,
			)
		}

		utils.LogInfo(fmt.Sprintf("Updating to version %s...", color.CyanString(release.TagName)))
		utils.LogDebug(fmt.Sprintf("downloading from %s", asset.DownloadUrl))
		data, err := http.Get(asset.DownloadUrl)
		if err != nil {
			return err
		}
		defer data.Body.Close()
		executablePath, err := os.Executable()
		if err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("current executable path as %s", executablePath))
		tempFile, err := utils.CreateTempFile(executablePath)
		if err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("created %s", tempFile.Name()))
		defer tempFile.Close()
		_, err = io.Copy(tempFile, data.Body)
		if err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("removing %s", executablePath))
		if err = os.Remove(executablePath); err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("renaming %s to %s", tempFile.Name(), executablePath))
		if err = os.Rename(tempFile.Name(), executablePath); err != nil {
			return err
		}
		utils.LogDebug(fmt.Sprintf("changing permissions of %s", executablePath))
		if err = os.Chmod(executablePath, 0755); err != nil {
			return err
		}
		utils.LogInfo(
			fmt.Sprintf(
				"%s Updated to version %s successfully!",
				utils.LogTickPrefix,
				color.CyanString(release.TagName),
			),
		)

		return nil
	},
}

func needsSelfUpdate() bool {
	release, err := core.GithubApiFetchLatestRelease(core.AppGithubOwner, core.AppGithubRepo)
	if err != nil {
		return false
	}
	return release.TagName != fmt.Sprintf("v%s", core.AppVersion)
}
