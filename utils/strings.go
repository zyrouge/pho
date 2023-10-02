package utils

import (
	"regexp"
	"strings"
)

var cleanTextInvalidRegex = regexp.MustCompile(`[^A-Za-z0-9-]+`)

func CleanId(text string) string {
	text = strings.ToLower(text)
	return cleanTextInvalidRegex.ReplaceAllLiteralString(text, "")
}
