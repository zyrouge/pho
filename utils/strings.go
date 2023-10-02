package utils

import (
	"regexp"
	"strings"
)

var cleanTextInvalidRegex = regexp.MustCompile(`[^A-Za-z0-9-]+`)

func ReplaceIllegalChars(text string) string {
	text = strings.ToLower(text)
	return cleanTextInvalidRegex.ReplaceAllLiteralString(text, "")
}

func CleanId(text string) string {
	text = ReplaceIllegalChars(text)
	text = strings.ToLower(text)
	return text
}

func CleanText(text string) string {
	text = ReplaceIllegalChars(text)
	return text
}
