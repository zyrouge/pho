package utils

import (
	"fmt"
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
