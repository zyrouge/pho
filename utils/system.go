package utils

import (
	"os/exec"
	"strings"
)

var ArchMap = map[string][]string{
	"amd64": {"amd64", "x86_64", "x86-64"},
	"386":   {"i386", "i686"},
	"arm64": {"armhf", "aarch64"},
	"arm":   {"arm"},
}

func GetSystemArch() string {
	raw, _ := ExecUnameM()
	for arch, aliases := range ArchMap {
		for _, alias := range aliases {
			if alias == raw {
				return arch
			}
		}
	}
	return ""
}

func ExecUnameM() (string, error) {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
