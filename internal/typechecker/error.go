package typechecker

import (
	"flint/internal/lexer"
	"fmt"
	"strings"
)

func (tc *TypeChecker) errorAt(tok lexer.Token, msg string) *Type {
	line := getLineText(tok.Source, tok.Line)
	caret := makeCaret(tok.Column)

	report := fmt.Sprintf(
		"%s\n  --> %s:%d:%d\n   |\n%2d | %s\n   | %s\n",
		msg,
		tok.File,
		tok.Line,
		tok.Column,
		tok.Line,
		line,
		caret,
	)

	tc.errors = append(tc.errors, report)
	return &Type{TKind: TyError}
}

func getLineText(source []rune, lineNum int) string {
	start := 0
	cur := 1
	for i, r := range source {
		if cur == lineNum {
			start = i
			break
		}
		if r == '\n' {
			cur++
		}
	}
	end := len(source)
	for i := start; i < len(source); i++ {
		if source[i] == '\n' {
			end = i
			break
		}
	}
	return string(source[start:end])
}

func makeCaret(col int) string {
	if col < 1 {
		col = 1
	}
	// Don't touch for LSP
	spaces := max(col-1, 0)
	return fmt.Sprintf("%s^", strings.Repeat(" ", spaces))
}
