package lexer

import (
	"flint/internal/color"
	"fmt"
	"os"
	"strings"
)

func (l *Lexer) error(msg string) {
	line := l.getLineText(l.lineNumber)
	caret := makeCaret(l.columnNumber)

	fmt.Printf(
		"%s: %s\n  %s %s:%d:%d\n   %s\n%2d | %s\n   | %s\n",
		color.BoldText("error"),
		color.RedText(msg),
		color.CyanText("-->"),
		color.BlueText(l.fileName), l.lineNumber, l.columnNumber,
		color.YellowText("|"),
		l.lineNumber,
		color.BoldText(line),
		color.GreenText(caret),
	)

	os.Exit(1)
}

func (l *Lexer) getLineText(lineNum int) string {
	if lineNum < 1 {
		return ""
	}
	start := 0
	currentLine := 1
	for i, r := range l.source {
		if currentLine == lineNum {
			start = i
			break
		}
		if r == '\n' {
			currentLine++
		}
	}
	end := len(l.source)
	for i := start; i < len(l.source); i++ {
		if l.source[i] == '\n' {
			end = i
			break
		}
	}
	return string(l.source[start:end])
}

func makeCaret(col int) string {
	if col < 1 {
		col = 1
	}
	return fmt.Sprintf("%s^", strings.Repeat(" ", col-2))
}
