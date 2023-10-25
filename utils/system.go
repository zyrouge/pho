package utils

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
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

type StartDetachedProcessOptions struct {
	Dir  string
	Exec string
	Args []string
}

func StartDetachedProcess(options *StartDetachedProcessOptions) error {
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
		Dir: options.Dir,
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
	proc, err := os.StartProcess(options.Exec, options.Args, procAttr)
	if err != nil {
		return err
	}
	if err = proc.Release(); err != nil {
		return err
	}
	return nil
}
