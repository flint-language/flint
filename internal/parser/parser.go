package parser

import (
	"flint/internal/lexer"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
	errors []string
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0, errors: []string{}}
}

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
	msg := fmt.Sprintf("expected token %v, got %v at %d:%d", kind, p.cur().Kind, p.cur().Line, p.cur().Column)
	p.error(msg)
	return p.cur(), false
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() (*Program, []string) {
	out := &Program{Exprs: []Expr{}}
	for p.cur().Kind != lexer.EndOfFile {
		if p.cur().Kind == lexer.Comment {
			p.eat()
			continue
		}
		expr := p.parseExpression(0)
		if expr == nil {
			p.eat()
			continue
		}
		out.Exprs = append(out.Exprs, expr)
	}
	return out, p.errors
}

func (p *Parser) parseExpression(minPrec int) Expr {
	if p.cur().Kind == lexer.KwVal {
		return p.parseValDecl()
	}
	if p.cur().Kind == lexer.KwMut {
		return p.parseMutDecl()
	}
	left := p.parsePrimary()
	if left == nil {
		return nil
	}
	for {
		opTok := p.cur()
		if opTok.Kind == lexer.Colon {
			p.eat()
			rightTok, ok := p.expect(lexer.Identifier)
			if !ok {
				return nil
			}
			left = &QualifiedExpr{
				Left:  left,
				Right: rightTok,
				Pos:   opTok,
			}
			continue
		}
		if opTok.Kind == lexer.Dot {
			p.eat()
			rightTok, ok := p.expect(lexer.Identifier)
			if !ok {
				return nil
			}
			left = &QualifiedExpr{
				Left:  left,
				Right: rightTok,
				Pos:   opTok,
			}
			continue
		}
		if opTok.Kind == lexer.EndOfFile {
			break
		}
		prec := opTok.Kind.Precedence()
		if prec == 0 || prec < minPrec {
			break
		}
		p.eat()
		nextMin := prec + 1
		right := p.parseExpression(nextMin)
		if right == nil {
			p.error(fmt.Sprintf("missing right-hand side after operator %q at %d:%d", opTok.Lexeme, opTok.Line, opTok.Column))
			return nil
		}
		if opTok.Kind == lexer.Pipe {
			left = &PipelineExpr{
				Left:  left,
				Right: right,
			}
		} else {
			left = &InfixExpr{
				Left:     left,
				Operator: opTok,
				Right:    right,
			}
		}
	}
	return left
}

func (p *Parser) parsePrimary() Expr {
	tok := p.cur()
	switch tok.Kind {
	case lexer.Identifier:
		p.eat()
		var expr Expr = &Identifier{Name: tok.Lexeme, Pos: tok}
		for p.cur().Kind == lexer.Colon || p.cur().Kind == lexer.Dot {
			opTok := p.eat()
			rightTok, ok := p.expect(lexer.Identifier)
			if !ok {
				return nil
			}
			expr = &QualifiedExpr{
				Left:  expr,
				Right: rightTok,
				Pos:   opTok,
			}
		}
		if p.cur().Kind == lexer.LeftParen {
			return p.parseCall(expr)
		}
		return expr
	case lexer.Int:
		p.eat()
		clean := stripNumericSeparators(tok.Lexeme)
		v, err := strconv.ParseInt(clean, 10, 64)
		if err != nil {
			p.error(fmt.Sprintf("invalid int literal %q at %d:%d", tok.Lexeme, tok.Line, tok.Column))
			return nil
		}
		return &IntLiteral{Value: v, Raw: tok.Lexeme, Pos: tok}
	case lexer.Float:
		p.eat()
		clean := stripNumericSeparators(tok.Lexeme)
		f, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			p.error(fmt.Sprintf("invalid float literal %q at %d:%d", tok.Lexeme, tok.Line, tok.Column))
			return nil
		}
		return &FloatLiteral{Value: f, Raw: tok.Lexeme, Pos: tok}
	case lexer.String:
		p.eat()
		return &StringLiteral{Value: tok.Lexeme, Pos: tok}
	case lexer.Byte:
		p.eat()
		if len(tok.Lexeme) != 3 || tok.Lexeme[0] != '\'' || tok.Lexeme[2] != '\'' {
			p.error(fmt.Sprintf("invalid byte literal %q at %d:%d", tok.Lexeme, tok.Line, tok.Column))
			return nil
		}
		return &ByteLiteral{
			Value: tok.Lexeme[1],
			Raw:   tok.Lexeme,
			Pos:   tok,
		}
	case lexer.Bool:
		p.eat()
		val := false
		if tok.Lexeme == "True" {
			val = true
		}
		return &BoolLiteral{Value: val, Pos: tok}
	case lexer.LeftParen:
		p.eat()
		if p.cur().Kind == lexer.RightParen {
			p.eat()
			return &TupleExpr{Elements: []Expr{}, Pos: tok}
		}
		elements := []Expr{}
		for {
			elem := p.parseExpression(0)
			if elem == nil {
				return nil
			}
			elements = append(elements, elem)
			if p.cur().Kind == lexer.Comma {
				p.eat()
				continue
			}
			break
		}
		if _, ok := p.expect(lexer.RightParen); !ok {
			return nil
		}
		if len(elements) == 1 {
			return elements[0]
		}
		return &TupleExpr{Elements: elements, Pos: tok}
	case lexer.Bang, lexer.Minus:
		p.eat()
		right := p.parseExpression(7)
		if right == nil {
			p.error(fmt.Sprintf("missing expression after prefix %q at %d:%d", tok.Lexeme, tok.Line, tok.Column))
			return nil
		}
		return &PrefixExpr{Operator: tok, Right: right}
	case lexer.KwPub:
		p.eat()
		switch p.cur().Kind {
		case lexer.KwFn:
			return p.parseFunc(true)
		case lexer.KwType:
			return p.parseTypeDecl(true)
		default:
			p.error("expected `fn` or `type` after `pub`")
			return nil
		}
	case lexer.KwFn:
		return p.parseFunc(false)
	case lexer.LeftBrace:
		return p.parseBlock()
	case lexer.KwUse:
		return p.parseUse()
	case lexer.KwIf:
		return p.parseIf()
	case lexer.KwMatch:
		return p.parseMatch()
	case lexer.KwFor:
		return p.parseForExpr()
	case lexer.LeftBracket:
		return p.parseList()
	case lexer.KwType:
		return p.parseTypeDecl(false)
	default:
		p.error(fmt.Sprintf("unexpected token %q (%v) at %d:%d", tok.Lexeme, tok.Kind, tok.Line, tok.Column))
		return nil
	}
}

func (p *Parser) parseCall(callee Expr) Expr {
	lparen := p.eat()
	args := []Expr{}
	if p.cur().Kind != lexer.RightParen {
		for {
			arg := p.parseExpression(0)
			if arg == nil {
				p.error("invalid argument expression in call")
				return nil
			}
			args = append(args, arg)
			if p.cur().Kind == lexer.Comma {
				p.eat()
				continue
			}
			break
		}
	}
	_, ok := p.expect(lexer.RightParen)
	if !ok {
		return nil
	}
	return &CallExpr{Callee: callee, Args: args, Pos: lparen}
}

func stripNumericSeparators(s string) string {
	out := []rune{}
	for _, r := range s {
		if r != '_' {
			out = append(out, r)
		}
	}
	return string(out)
}

func (p *Parser) parseValDecl() Expr {
	p.eat()
	nameTok, ok := p.expect(lexer.Identifier)
	if !ok {
		return nil
	}
	var typeAnn Expr
	if p.cur().Kind == lexer.Colon {
		p.eat()
		typeAnn = p.parseType()
	}
	_, ok = p.expect(lexer.Equal)
	if !ok {
		return nil
	}
	value := p.parseExpression(0)
	if value == nil {
		p.error(fmt.Sprintf("missing initializer for val %s", nameTok.Lexeme))
		return nil
	}
	return &ValDeclExpr{
		Name:  nameTok,
		Type:  typeAnn,
		Value: value,
	}
}

func (p *Parser) parseMutDecl() Expr {
	p.eat()
	nameTok, ok := p.expect(lexer.Identifier)
	if !ok {
		return nil
	}
	var typeAnn Expr
	if p.cur().Kind == lexer.Colon {
		p.eat()
		typeAnn = p.parseType()
	}
	_, ok = p.expect(lexer.Equal)
	if !ok {
		return nil
	}
	value := p.parseExpression(0)
	if value == nil {
		p.error(fmt.Sprintf("missing initializer for mut %s", nameTok.Lexeme))
		return nil
	}
	return &MutDeclExpr{
		Name:  nameTok,
		Type:  typeAnn,
		Value: value,
	}
}

func (p *Parser) parseFunc(pub bool) Expr {
	p.expect(lexer.KwFn)
	nameTok, ok := p.expect(lexer.Identifier)
	if !ok {
		return nil
	}
	p.expect(lexer.LeftParen)
	params := []Param{}
	if p.cur().Kind != lexer.RightParen {
		for {
			paramTok, ok := p.expect(lexer.Identifier)
			if !ok {
				return nil
			}

			var typ Expr
			if p.cur().Kind == lexer.Colon {
				p.eat()
				typ = p.parseType()
			}
			params = append(params, Param{Name: paramTok, Type: typ})
			if p.cur().Kind == lexer.Comma {
				p.eat()
				continue
			}
			break
		}
	}
	p.expect(lexer.RightParen)
	var retType Expr
	if p.cur().Kind != lexer.LeftBrace {
		retType = p.parseType()
	}
	body := p.parseBlock()
	if body == nil {
		return nil
	}
	return &FuncDeclExpr{
		Pub:    pub,
		Name:   nameTok,
		Params: params,
		Ret:    retType,
		Body:   body,
	}
}

func (p *Parser) parseBlock() Expr {
	p.expect(lexer.LeftBrace)
	exprs := []Expr{}
	for p.cur().Kind != lexer.RightBrace && p.cur().Kind != lexer.EndOfFile {
		if p.cur().Kind == lexer.Comment {
			p.eat()
			continue
		}
		e := p.parseExpression(0)
		if e == nil {
			p.eat()
			continue
		}
		exprs = append(exprs, e)
	}
	p.expect(lexer.RightBrace)
	return &BlockExpr{Exprs: exprs}
}

func (p *Parser) parseUse() Expr {
	start := p.eat()
	path := []string{}
	for {
		tok, ok := p.expect(lexer.Identifier)
		if !ok {
			return nil
		}
		path = append(path, tok.Lexeme)

		if p.cur().Kind == lexer.Slash {
			p.eat()
			continue
		}
		break
	}
	members := []string{}
	alias := ""
	if p.cur().Kind == lexer.Dot && p.peek(1).Kind == lexer.LeftBrace {
		p.eat()
		p.eat()
		for {
			memberTok, ok := p.expect(lexer.Identifier)
			if !ok {
				return nil
			}
			members = append(members, memberTok.Lexeme)

			if p.cur().Kind == lexer.Comma {
				p.eat()
				continue
			}
			break
		}
		p.expect(lexer.RightBrace)
	}
	if p.cur().Kind == lexer.KwAs {
		p.eat()
		aliasTok, ok := p.expect(lexer.Identifier)
		if !ok {
			return nil
		}
		alias = aliasTok.Lexeme
	}
	return &UseExpr{
		Path:    path,
		Alias:   alias,
		Members: members,
		Pos:     start,
	}
}

func (p *Parser) parseIf() Expr {
	p.eat()
	cond := p.parseExpression(0)
	if cond == nil {
		p.error("expected condition after 'if'")
		return nil
	}
	thenBlock := p.parseBlock()
	if thenBlock == nil {
		p.error("expected block after 'if' condition")
		return nil
	}
	var elseBlock Expr
	if p.cur().Kind == lexer.KwElse {
		p.eat()
		elseBlock = p.parseBlock()
		if elseBlock == nil {
			p.error("expected block after 'else'")
			return nil
		}
	}
	return &IfExpr{
		Cond: cond,
		Then: thenBlock,
		Else: elseBlock,
	}
}

func (p *Parser) parseMatch() Expr {
	p.eat()
	value := p.parseExpression(0)
	if value == nil {
		p.error("expected expression after 'match'")
		return nil
	}
	_, ok := p.expect(lexer.LeftBrace)
	if !ok {
		return nil
	}
	arms := []*MatchArm{}
	for p.cur().Kind != lexer.RightBrace && p.cur().Kind != lexer.EndOfFile {
		if p.cur().Kind == lexer.Vbar {
			p.eat()
		}
		var pattern Expr
		if p.cur().Kind == lexer.Underscore {
			pattern = &Identifier{Name: "_", Pos: p.cur()}
			p.eat()
		} else {
			pattern = p.parseExpression(0)
			if pattern == nil {
				p.error("expected pattern in match arm")
				return nil
			}
		}
		var guard Expr
		if p.cur().Kind == lexer.KwIf {
			p.eat()
			guard = p.parseExpression(0)
			if guard == nil {
				p.error("expected guard expression after 'if'")
				return nil
			}
		}
		_, ok := p.expect(lexer.RArrow)
		if !ok {
			return nil
		}
		body := p.parseExpression(0)
		if body == nil {
			p.error("expected body expression in match arm")
			return nil
		}
		arms = append(arms, &MatchArm{
			Pattern: pattern,
			Guard:   guard,
			Body:    body,
		})
	}
	_, ok = p.expect(lexer.RightBrace)
	if !ok {
		return nil
	}
	return &MatchExpr{
		Value: value,
		Arms:  arms,
	}
}

func (p *Parser) parseForExpr() Expr {
	p.eat()
	var vars []Expr
	for {
		expr := p.parsePrimary()
		if expr == nil {
			return nil
		}
		vars = append(vars, expr)

		if p.cur().Kind == lexer.Comma {
			p.eat()
			continue
		}
		break
	}
	_, ok := p.expect(lexer.KwIn)
	if !ok {
		return nil
	}
	iterable := p.parseExpression(0)
	if iterable == nil {
		return nil
	}
	var whereExpr Expr
	if p.cur().Kind == lexer.KwWhere {
		p.eat()
		whereExpr = p.parseExpression(0)
	}
	body := p.parseBlock()
	if body == nil {
		return nil
	}
	return &ForExpr{
		Vars:     vars,
		Iterable: iterable,
		Where:    whereExpr,
		Body:     body,
	}
}

func (p *Parser) parseList() Expr {
	start := p.eat()
	elements := []Expr{}
	for p.cur().Kind != lexer.RightBracket && p.cur().Kind != lexer.EndOfFile {
		elem := p.parseExpression(0)
		if elem == nil {
			p.error("invalid list element")
			return nil
		}
		elements = append(elements, elem)
		if p.cur().Kind == lexer.Comma {
			p.eat()
			continue
		}
		break
	}
	if _, ok := p.expect(lexer.RightBracket); !ok {
		return nil
	}
	return &ListExpr{
		Elements: elements,
		Pos:      start,
	}
}

func (p *Parser) parseType() Expr {
	tok := p.cur()
	switch tok.Kind {
	case lexer.KwInt, lexer.KwFloat, lexer.KwBool, lexer.KwByte, lexer.KwString, lexer.KwNil:
		p.eat()
		return &TypeExpr{Name: tok.Lexeme, Pos: tok}
	case lexer.KwList:
		p.eat()
		var elemType Expr
		if p.cur().Kind == lexer.LeftParen {
			p.eat()
			elemType = p.parseType()
			p.expect(lexer.RightParen)
		}
		return &TypeExpr{Name: "List", Pos: tok, Generic: elemType}
	case lexer.LeftParen:
		p.eat()
		types := []Expr{}
		for {
			t := p.parseType()
			if t == nil {
				return nil
			}
			types = append(types, t)
			if p.cur().Kind == lexer.Comma {
				p.eat()
				continue
			}
			break
		}
		p.expect(lexer.RightParen)
		return &TupleTypeExpr{Types: types, Pos: tok}
	default:
		p.error(fmt.Sprintf("expected type, got %q (%v) at %d:%d", tok.Lexeme, tok.Kind, tok.Line, tok.Column))
		return nil
	}
}

func (p *Parser) parseTypeDecl(pub bool) Expr {
	p.expect(lexer.KwType)
	nameTok, ok := p.expect(lexer.Identifier)
	if !ok {
		return nil
	}
	var body Expr
	if p.cur().Kind == lexer.LeftBrace {
		p.eat()
		fields := []Param{}
		for p.cur().Kind != lexer.RightBrace && p.cur().Kind != lexer.EndOfFile {
			if p.cur().Kind == lexer.Comment {
				p.eat()
				continue
			}
			fieldTok, ok := p.expect(lexer.Identifier)
			if !ok {
				return nil
			}
			if p.cur().Kind != lexer.Colon {
				p.error(fmt.Sprintf("expected ':' after field name %s", fieldTok.Lexeme))
				return nil
			}
			p.eat()
			fieldType := p.parseType()
			if fieldType == nil {
				return nil
			}
			fields = append(fields, Param{Name: fieldTok, Type: fieldType})
			if p.cur().Kind == lexer.Comma {
				p.eat()
			}
		}
		p.expect(lexer.RightBrace)
		body = &RecordTypeExpr{Name: nameTok, Fields: fields, Pos: nameTok}
	}
	return &TypeDeclExpr{
		Pub:  pub,
		Name: nameTok,
		Body: body,
		Pos:  nameTok,
	}
}
