package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/zyrouge/pho/utils"
)

type DeflatedAppImage struct {
	AppImagePath string
	ParentDir    string
	AppDir       string
}

func DeflateAppImage(appImagePath string, parentDir string) (*DeflatedAppImage, error) {
	cmd := exec.Command(appImagePath, "--appimage-extract")
	cmd.Dir = parentDir
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	appDir := path.Join(parentDir, "squashfs-root")
	deflated := &DeflatedAppImage{
		AppImagePath: appImagePath,
		ParentDir:    parentDir,
		AppDir:       appDir,
	}
	return deflated, nil
}

type DeflatedAppImageMetadata struct {
	*DeflatedAppImage
	ExecName    string
	IconPath    string
	DesktopPath string
}

func (deflated *DeflatedAppImage) ExtractMetadata() (*DeflatedAppImageMetadata, error) {
	execName, err := deflated.ExtractExecName()
	if err != nil {
		return nil, err
	}
	desktopPath := path.Join(deflated.AppDir, fmt.Sprintf("%s.desktop", execName))
	iconPath := path.Join(deflated.AppDir, ".DirIcon")
	metadata := &DeflatedAppImageMetadata{
		DeflatedAppImage: deflated,
		ExecName:         execName,
		IconPath:         iconPath,
		DesktopPath:      desktopPath,
	}
	return metadata, nil
}

func (deflated *DeflatedAppImage) ExtractExecName() (string, error) {
	files, err := os.ReadDir(deflated.AppDir)
	if err != nil {
		return "", err
	}
	for _, x := range files {
		name := x.Name()
		if strings.HasSuffix(name, ".desktop") {
			return strings.TrimSuffix(name, ".desktop"), nil
		}
	}
	return "", errors.New("cannot find .desktop file from AppDir")
}

func (metadata *DeflatedAppImageMetadata) CopyIconFile(paths *AppPaths) error {
	src, err := os.Open(metadata.IconPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dest, err := os.Create(paths.Icon)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return err
}

var desktopFileExecRegex = regexp.MustCompile(`Exec=("[^"]+"|[^\s]+)`)
var desktopFileIconRegex = regexp.MustCompile(`Icon=[^\n]+`)

func (metadata *DeflatedAppImageMetadata) InstallDesktopFile(paths *AppPaths) error {
	bytes, err := os.ReadFile(metadata.DesktopPath)
	if err != nil {
		return err
	}
	return InstallDesktopFile(paths, string(bytes))
}

func InstallDesktopFile(paths *AppPaths, content string) error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}
	execPath := utils.QuotedWhenSpace(paths.AppImage)
	if !config.EnableIntegrationPrompt {
		execPath = "env APPIMAGELAUNCHER_DISABLE=1 " + execPath
	}
	content = replaceDesktopEntry(
		content,
		desktopFileExecRegex,
		fmt.Sprintf("Exec=%s", utils.QuotedWhenSpace(execPath)),
	)
	content = replaceDesktopEntry(
		content,
		desktopFileIconRegex,
		fmt.Sprintf("Icon=%s", utils.QuotedWhenSpace(paths.Icon)),
	)
	content = strings.TrimSpace(content)
	if err := os.WriteFile(paths.Desktop, []byte(content), os.ModePerm); err != nil {
		return err
	}
	cmd := exec.Command("xdg-desktop-menu", "install", paths.Desktop, "--novendor")
	return cmd.Run()
}

func UninstallDesktopFile(desktopFilePath string) error {
	cmd := exec.Command("xdg-desktop-menu", "uninstall", desktopFilePath, "--novendor")
	return cmd.Run()
}

func (metadata *DeflatedAppImageMetadata) Symlink(paths *AppPaths) error {
	if err := os.Symlink(paths.AppImage, paths.Symlink); err != nil {
		return err
	}
	return nil
}

func replaceDesktopEntry(content string, pattern *regexp.Regexp, replaceWith string) string {
	count := 0
	content = pattern.ReplaceAllStringFunc(content, func(s string) string {
		count++
		return replaceWith
	})
	if count == 0 {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += replaceWith
	}
	return content
}
