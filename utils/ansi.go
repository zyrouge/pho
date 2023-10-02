package utils

import (
	"fmt"
	"regexp"
)

// Source: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797

const AnsiEsc = "\x1B"
const AnsiResetCursor = "\r"

var AnsiEraseLine = fmt.Sprintf("%s[K", AnsiEsc)

func AnsiCursorLineUp(lines int) string {
	return fmt.Sprintf("%s[%dA", AnsiEsc, lines)
}

func AnsiCursorToColumn(offset int) string {
	return fmt.Sprintf("%s[%dG", AnsiEsc, offset)
}

// Source: https://stackoverflow.com/questions/17998978/removing-colors-from-output#comment117315951_35582778
var stripAnsiRegex = regexp.MustCompile(`\x1B\[(?:;?[0-9]{1,3})+[mGK]`)

func StripAnsi(text string) string {
	return stripAnsiRegex.ReplaceAllLiteralString(text, "")
}
