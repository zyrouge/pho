package utils

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ReaderReadLine(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func ResolvePath(name string) (string, error) {
	if name == "" {
		return "", errors.New("cannot resolve to empty path")
	}
	if strings.HasPrefix(name, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		name = filepath.Join(home, name[2:])
	}
	name, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	return name, err
}
