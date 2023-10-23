package utils

import (
	"encoding/json"
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

var titleCaseReplaceRegex = regexp.MustCompile(`\w{2,}`)

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

func PrettyJSONString(data any) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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
