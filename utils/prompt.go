package utils

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PromptTextInput(reader *bufio.Reader, question string, defaultValue string) (string, error) {
	prefix := fmt.Sprintf(
		"%s %s",
		LogQuestionPrefix,
		question,
	)
	fmt.Printf("%s %s ", prefix, color.HiBlackString(defaultValue))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		input = defaultValue
	}
	fmt.Print(AnsiCursorLineUp(1), AnsiEraseLine, AnsiResetCursor)
	fmt.Printf("%s %s\n", prefix, color.CyanString(input))
	return input, nil
}

func PromptYesNoInput(reader *bufio.Reader, question string) (bool, error) {
	prefix := fmt.Sprintf(
		"%s %s",
		LogQuestionPrefix,
		question,
	)
	fmt.Printf("%s %s ", prefix, color.HiBlackString("[y/N]"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = strings.TrimSuffix(input, "\n")
	input = strings.ToLower(input)
	if input == "y" || input == "yes" {
		input = "yes"
	} else {
		input = "no"
	}
	TerminalErasePreviousLine()
	fmt.Printf("%s %s\n", prefix, color.CyanString(input))
	return input == "yes", nil
}
