package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var LogDebugPrefix = color.HiBlackString("debug:")
var LogWarnPrefix = color.YellowString("warning")
var LogErrorPrefix = color.RedString("error:")

var LogQuestionPrefix = color.GreenString("?")
var LogRightArrowPrefix = color.MagentaString(">")
var LogTickPrefix = color.GreenString("√")
var LogExclamationPrefix = color.RedString("!")

var LogDebugEnabled = false

func LogInfo(msg string) {
	fmt.Println(msg)
}

func LogDebug(msg string) {
	if !LogDebugEnabled {
		return
	}
	fmt.Printf("%s %s\n", LogDebugPrefix, msg)
}

func LogWarning(msg string) {
	fmt.Printf("%s %s\n", LogWarnPrefix, msg)
}

func LogError(err any) {
	fmt.Printf("%s %v\n", LogErrorPrefix, err)
}

func LogLn() {
	fmt.Println()
}

type LogTable struct {
	Columns         [][]string
	RowsCount       int
	ColumnsCount    int
	RowWidths       []int
	RowSeparator    string
	ColumnSeparator string
}

func NewLogTable() *LogTable {
	return &LogTable{
		Columns:         [][]string{},
		RowsCount:       0,
		ColumnsCount:    0,
		RowWidths:       []int{},
		RowSeparator:    "\n",
		ColumnSeparator: "   ",
	}
}

func (table *LogTable) Add(values ...string) {
	table.AddColumn(values)
}

func (table *LogTable) AddColumn(column []string) {
	table.Columns = append(table.Columns, column)
	table.RowsCount++
	columnCount := len(column)
	if columnCount > table.ColumnsCount {
		table.RowWidths = append(
			table.RowWidths,
			make([]int, columnCount-table.ColumnsCount)...,
		)
		table.ColumnsCount = columnCount
	}
	for i, x := range column {
		currWidth := len(StripAnsi(x))
		if currWidth > table.RowWidths[i] {
			table.RowWidths[i] = currWidth
		}
	}
}

func (table *LogTable) Build() string {
	var text strings.Builder
	for i, row := range table.Columns {
		isLastRow := i == table.RowsCount-1
		for j, x := range row {
			isLastColumn := j == table.ColumnsCount-1
			maxWidth := table.RowWidths[j]
			currWidth := len(StripAnsi(x))
			text.WriteString(x)
			text.WriteString(strings.Repeat(" ", maxWidth-currWidth))
			if !isLastColumn {
				text.WriteString(table.ColumnSeparator)
			}
		}
		if !isLastRow {
			text.WriteString(table.RowSeparator)
		}
	}
	return text.String()
}

func (table *LogTable) Print() {
	fmt.Println(table.Build())
}
