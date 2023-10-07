package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var cleanTextInvalidRegex = regexp.MustCompile(`[^A-Za-z0-9-]+`)
var cleanTextTrimRegex = regexp.MustCompile(`(^-|-$)`)

func ReplaceIllegalChars(text string) string {
	text = cleanTextInvalidRegex.ReplaceAllLiteralString(text, "")
	text = cleanTextTrimRegex.ReplaceAllLiteralString(text, "")
	return text
}

func CleanText(text string) string {
	text = ReplaceIllegalChars(text)
	return text
}

func CleanId(text string) string {
	text = ReplaceIllegalChars(text)
	text = strings.ToLower(text)
	return text
}

func CleanName(text string) string {
	text = ReplaceIllegalChars(text)
	text = TitleCase(text)
	return text
}

var titleCaseReplaceRegex = regexp.MustCompile(`/(\w)(\w+)/`)

func TitleCase(text string) string {
	return titleCaseReplaceRegex.ReplaceAllStringFunc(
		text,
		func(s string) string {
			runes := []rune(s)
			runes[0] = unicode.ToUpper(runes[0])
			return string(runes)
		},
	)
}
