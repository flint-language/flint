package parser

import (
	"flint/internal/lexer"
	"fmt"
	"strings"
)

func (p *Parser) errorAt(tok lexer.Token, msg string) {
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

	p.errors = append(p.errors, report)
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

func (p *Parser) synchronize() {
	for p.cur().Kind != lexer.EndOfFile {
		switch p.cur().Kind {
		case lexer.KwFn, lexer.KwVal, lexer.KwMut, lexer.KwIf,
			lexer.KwType, lexer.KwMatch, lexer.KwUse:
			return
		case lexer.RightBrace:
			p.eat()
			return
		default:
			p.eat()
		}
	}
}
