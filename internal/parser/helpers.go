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

func (p *Parser) coerceExprToType(expr Expr, typ Expr) Expr {
	switch e := expr.(type) {
	case *IntLiteral:
		if t, ok := typ.(*TypeExpr); ok {
			switch t.Name {
			case "U32", "U64", "U16", "U8", "Unsigned":
				return &UnsignedLiteral{Value: uint64(e.Value), Raw: e.Raw, Pos: e.Pos}
			}
		}
	case *UnsignedLiteral:
		if t, ok := typ.(*TypeExpr); ok {
			switch t.Name {
			case "I32", "I64", "I16", "I8", "Int":
				return &IntLiteral{Value: int64(e.Value), Raw: e.Raw, Pos: e.Pos}
			}
		}
	case *FloatLiteral:
		if t, ok := typ.(*TypeExpr); ok {
			switch t.Name {
			case "F32", "F64", "Float":
				return &FloatLiteral{Value: e.Value, Raw: e.Raw, Pos: e.Pos}
			}
		}
	case *InfixExpr:
		e.Left = p.coerceExprToType(e.Left, typ)
		e.Right = p.coerceExprToType(e.Right, typ)
	case *PrefixExpr:
		e.Right = p.coerceExprToType(e.Right, typ)
	case *CallExpr:
		for i := range e.Args {
			e.Args[i] = p.coerceExprToType(e.Args[i], typ)
		}
	case *BlockExpr:
		for i := range e.Exprs {
			e.Exprs[i] = p.coerceExprToType(e.Exprs[i], typ)
		}
	case *TupleExpr:
		for i := range e.Elements {
			e.Elements[i] = p.coerceExprToType(e.Elements[i], typ)
		}
	case *ListExpr:
		for i := range e.Elements {
			e.Elements[i] = p.coerceExprToType(e.Elements[i], typ)
		}
	}
	return expr
}
