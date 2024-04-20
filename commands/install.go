package commands

import (
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"github.com/zyrouge/pho/core"
	"github.com/zyrouge/pho/utils"
)

var InstallCommand = cli.Command{
	Name:    "install",
	Aliases: []string{"add"},
	Usage:   "Install an application",
	Commands: []*cli.Command{
		&InstallGithubCommand,
		&InstallLocalCommand,
		&InstallHttpCommand,
	},
}

type Installable struct {
	Name        string
	Id          string
	DownloadUrl string
	Size        int
}

type InstallableAppStatus int

const (
	InstallableAppFailed InstallableAppStatus = iota
	InstallableAppDownloading
	InstallableAppIntegrating
	InstallableAppInstalled
)

type InstallableApp struct {
	App    *core.AppConfig
	Source any
	Asset  *core.Asset

	Index          int
	Count          int
	StartedAt      int64
	Progress       int64
	RawProgress    InstallableAppRawProgress
	Speed          int64
	RemainingSecs  int64
	PrintCycle     int
	SkipCycleErase bool
	Status         InstallableAppStatus
}

type InstallableAppRawProgress struct {
	Sizes []int
	Times []int64
}

func (x *InstallableApp) Write(data []byte) (n int, err error) {
	l := len(data)
	x.Progress += int64(l)
	x.RawProgress.Sizes = append(x.RawProgress.Sizes, l)
	x.RawProgress.Times = append(x.RawProgress.Times, utils.TimeNowSeconds())
	remove := len(x.RawProgress.Sizes) - 50
	if remove > 0 {
		x.RawProgress.Sizes = slices.Delete(x.RawProgress.Sizes, 0, remove)
		x.RawProgress.Times = slices.Delete(x.RawProgress.Times, 0, remove)
	}
	return l, nil
}

func (x *InstallableApp) logDebug(msg string) {
	if utils.LogDebugEnabled {
		x.SkipCycleErase = true
		utils.LogDebug(msg)
	}
}

func (x *InstallableApp) PrintStatus() {
	if x.PrintCycle > 0 && !x.SkipCycleErase {
		utils.TerminalErasePreviousLine()
		x.SkipCycleErase = false
	}
	x.PrintCycle++

	prefix := color.HiBlackString(fmt.Sprintf("[%d/%d]", x.Index+1, x.Count))
	elapsedSecs := utils.HumanizeSeconds(utils.TimeNowSeconds() - x.StartedAt)
	suffix := color.HiBlackString(fmt.Sprintf("(%s)", elapsedSecs))

	switch x.Status {
	case InstallableAppFailed:
		fmt.Printf(
			"%s %s %s %s %s\n",
			prefix,
			utils.LogExclamationPrefix,
			color.RedString(x.App.Id),
			x.App.Version,
			suffix,
		)

	case InstallableAppDownloading:
		x.calculateMetrics()
		suffix := color.HiBlackString(
			fmt.Sprintf(
				"(%s / %s @ %s/s)",
				elapsedSecs,
				utils.HumanizeSeconds(x.RemainingSecs),
				prettyBytes(x.Speed),
			),
		)
		fmt.Printf(
			"%s %s %s %s (%s / %s) %s\n",
			prefix,
			color.YellowString(utils.TerminalLoadingSymbol(x.PrintCycle)),
			color.CyanString(x.App.Id),
			x.App.Version,
			prettyBytes(x.Progress),
			prettyBytes(x.Asset.Size),
			suffix,
		)

	case InstallableAppIntegrating:
		fmt.Printf(
			"%s %s %s %s %s\n",
			prefix,
			color.YellowString(utils.TerminalLoadingSymbol(x.PrintCycle)),
			color.CyanString(x.App.Id),
			x.App.Version,
			suffix,
		)

	case InstallableAppInstalled:
		fmt.Printf(
			"%s %s %s %s %s\n",
			prefix,
			utils.LogTickPrefix,
			color.GreenString(x.App.Id),
			x.App.Version,
			suffix,
		)
	}
}

const printStatusTickerDuration = time.Second / 4

func (x *InstallableApp) StartStatusTicker() *time.Ticker {
	ticker := time.NewTicker(printStatusTickerDuration)
	go func() {
		for range ticker.C {
			x.PrintStatus()
		}
	}()
	return ticker
}

func InstallApps(apps []InstallableApp) (int, int) {
	success := 0
	count := len(apps)
	for i := range apps {
		x := &apps[i]
		x.Index = i
		x.Count = count
		x.StartedAt = utils.TimeNowSeconds()
		x.Status = InstallableAppDownloading
		x.RawProgress = InstallableAppRawProgress{
			Sizes: []int{},
			Times: []int64{},
		}
		x.PrintStatus()
		x.logDebug("updating transactions")
		core.UpdateTransactions(func(transactions *core.Transactions) error {
			transactions.PendingInstallations[x.App.Id] = core.PendingInstallation{
				InvolvedDirs:  []string{x.App.Paths.Dir},
				InvolvedFiles: []string{x.App.Paths.Desktop},
			}
			return nil
		})
		if err := x.Install(); err != nil {
			x.Status = InstallableAppFailed
			x.PrintStatus()
			utils.LogError(err)
			break
		} else {
			x.Status = InstallableAppInstalled
			x.PrintStatus()
			success++
		}
		x.logDebug("updating transactions")
		core.UpdateTransactions(func(transactions *core.Transactions) error {
			delete(transactions.PendingInstallations, x.App.Id)
			return nil
		})
	}
	return success, count - success
}

func (x *InstallableApp) Install() error {
	ticker := x.StartStatusTicker()
	defer ticker.Stop()
	if err := x.Download(); err != nil {
		return err
	}
	x.Status = InstallableAppIntegrating
	if err := x.Integrate(); err != nil {
		return err
	}
	if err := x.SaveConfig(); err != nil {
		return err
	}
	return nil
}

func (x *InstallableApp) Download() error {
	x.logDebug(fmt.Sprintf("creating %s", x.App.Paths.Dir))
	if err := os.MkdirAll(x.App.Paths.Dir, os.ModePerm); err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("creating %s", x.App.Paths.Desktop))
	if err := os.MkdirAll(path.Dir(x.App.Paths.Desktop), os.ModePerm); err != nil {
		return err
	}
	tempFile, err := utils.CreateTempFile(x.App.Paths.AppImage)
	if err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("created %s", tempFile.Name()))
	defer tempFile.Close()
	data, err := x.Asset.Download()
	if err != nil {
		return err
	}
	defer data.Close()
	mw := io.MultiWriter(tempFile, x)
	_, err = io.Copy(mw, data)
	if err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("renaming %s to %s", tempFile.Name(), x.App.Paths.AppImage))
	if err = os.Rename(tempFile.Name(), x.App.Paths.AppImage); err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("changing permissions of %s", x.App.Paths.AppImage))
	return os.Chmod(x.App.Paths.AppImage, 0755)
}

func (x *InstallableApp) Integrate() error {
	tempDir := path.Join(x.App.Paths.Dir, "temp")
	x.logDebug(fmt.Sprintf("creating %s", tempDir))
	err := os.Mkdir(tempDir, os.ModePerm)
	if err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("deflating %s into %s", x.App.Paths.AppImage, tempDir))
	deflated, err := core.DeflateAppImage(x.App.Paths.AppImage, tempDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	metadata, err := deflated.ExtractMetadata()
	if err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("creating %s", x.App.Paths.Icon))
	if err = metadata.CopyIconFile(&x.App.Paths); err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("installing .desktop file at %s", x.App.Paths.Desktop))
	if err = metadata.InstallDesktopFile(&x.App.Paths); err != nil {
		return err
	}
	if x.App.Paths.Symlink != "" {
		x.logDebug(fmt.Sprintf("creating symlink %s", x.App.Paths.Symlink))
		if err = metadata.Symlink(&x.App.Paths); err != nil {
			return err
		}
	}
	return nil
}

func (x *InstallableApp) SaveConfig() error {
	x.logDebug(fmt.Sprintf("saving app config to %s", x.App.Paths.Config))
	if err := core.SaveAppConfig(x.App.Paths.Config, x.App); err != nil {
		return err
	}
	x.logDebug(fmt.Sprintf("saving app source config to %s", x.App.Paths.SourceConfig))
	if err := core.SaveSourceConfig[any](x.App.Paths.SourceConfig, x.Source); err != nil {
		return err
	}
	config, err := core.ReadConfig()
	if err != nil {
		return err
	}
	config.Installed[x.App.Id] = x.App.Paths.Config
	x.logDebug("saving config")
	return core.SaveConfig(config)
}

func prettyBytes(size int64) string {
	if size < 1000 {
		kb := float32(size) / 1000
		return fmt.Sprintf("%.2f KB", kb)
	}
	mb := float32(size) / 1000000
	return fmt.Sprintf("%.2f MB", mb)
}

func (x *InstallableApp) calculateMetrics() {
	count := len(x.RawProgress.Sizes)
	if count < 2 {
		return
	}
	total := 0
	for _, x := range x.RawProgress.Sizes {
		total += x
	}
	time := max(1, x.RawProgress.Times[count-1]-x.RawProgress.Times[0])
	x.Speed = int64(total) / time
	if x.Speed > 0 {
		x.RemainingSecs = (x.Asset.Size - x.Progress) / x.Speed
	} else {
		x.RemainingSecs = 0
	}
}
