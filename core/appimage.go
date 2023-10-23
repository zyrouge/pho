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
	err := cmd.Run()
	if err != nil {
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

var desktopFileExecRegex = regexp.MustCompile(`Exec=("[^"]+"|[^ ]+)`)
var desktopFileIconRegex = regexp.MustCompile(`Icon=[^\n]+`)

func (metadata *DeflatedAppImageMetadata) InstallDesktopFile(paths *AppPaths) error {
	bytes, err := os.ReadFile(metadata.DesktopPath)
	if err != nil {
		return err
	}
	content := string(bytes)
	content = desktopFileExecRegex.ReplaceAllLiteralString(
		content,
		fmt.Sprintf("Exec=%s", utils.QuotedWhenSpace(paths.AppImage)),
	)
	content = desktopFileIconRegex.ReplaceAllLiteralString(
		content,
		fmt.Sprintf("Icon=%s", utils.QuotedWhenSpace(paths.Icon)),
	)
	content = strings.TrimSpace(content)
	if err = os.WriteFile(paths.Desktop, []byte(content), os.ModePerm); err != nil {
		return err
	}
	cmd := exec.Command("xdg-desktop-menu", "install", paths.Desktop, "--novendor")
	return cmd.Run()
}

func UninstallDesktopFile(desktopFilePath string) error {
	cmd := exec.Command("xdg-desktop-menu", "uninstall", desktopFilePath, "--novendor")
	return cmd.Run()
}
