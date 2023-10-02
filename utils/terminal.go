package utils

import "fmt"

func TerminalErasePreviousLine() {
	fmt.Print(AnsiCursorLineUp(1), AnsiEraseLine, AnsiResetCursor)
}

// \ | / -
func TerminalLoadingSymbol(v int) string {
	v = v % 4
	switch v {
	case 1:
		return "\\"

	case 2:
		return "|"

	case 3:
		return "/"

	default:
		return "-"
	}
}
