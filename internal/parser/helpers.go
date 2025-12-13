package parser

import (
	"flint/internal/lexer"
	"fmt"
)

func (p *Parser) cur() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Kind: lexer.EndOfFile, Lexeme: "", Line: 0, Column: 0}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek(n int) lexer.Token {
	i := p.pos + n
	if i >= len(p.tokens) {
		return lexer.Token{Kind: lexer.EndOfFile, Lexeme: "", Line: 0, Column: 0}
	}
	return p.tokens[i]
}

func (p *Parser) eat() lexer.Token {
	t := p.cur()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return t
}

func (p *Parser) expect(kind lexer.TokenKind) (lexer.Token, bool) {
	if p.cur().Kind == kind {
		return p.eat(), true
	}
	tok := p.cur()
	p.errorAt(tok, fmt.Sprintf(
		"expected token %v, got %v",
		kind, tok.Kind,
	))
	p.synchronize()
	return tok, false
}

func (p *Parser) attachDocs(expr Expr) {
	if p.pendingDocs == nil {
		return
	}
	switch n := expr.(type) {
	case *FuncDeclExpr:
		n.Docs = p.pendingDocs
	}
	p.pendingDocs = nil
}
