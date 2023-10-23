package utils

import (
	"regexp"
	"strings"
)

var cleanTextInvalidRegex = regexp.MustCompile(`[^A-Za-z0-9-]+`)
var cleanTextTrimRegex = regexp.MustCompile(`(^-|-$)`)

func ReplaceIllegalChars(text string) string {
	text = cleanTextInvalidRegex.ReplaceAllLiteralString(text, "")
	text = cleanTextTrimRegex.ReplaceAllLiteralString(text, "")
	return text
}

func CleanId(text string) string {
	text = ReplaceIllegalChars(text)
	text = strings.ToLower(text)
	return text
}

func QuotedWhenSpace(text string) string {
	return EncloseWhen(text, " ", `"`, `"`)
}

func EncloseWhen(text string, when string, start string, end string) string {
	if !strings.Contains(text, when) {
		return text
	}
	return start + text + end
}
