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

func ParseProgram(tokens []lexer.Token) (*Program, []string) {
	p := new(tokens)
	out := &Program{Exprs: []Expr{}}
	for p.cur().Kind != lexer.EndOfFile {
		if p.cur().Kind == lexer.Comment {
			p.eat()
			continue
		}
		decorators := p.parseDecorators()
		expr := p.parseExpression(0)
		if expr == nil {
			p.synchronize()
			continue
		}
		if fn, ok := expr.(*FuncDeclExpr); ok {
			fn.Decorators = decorators
		}
		out.Exprs = append(out.Exprs, expr)
	}
	dectectRecursion(out)
	return out, p.errors
}

func new(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0, errors: []string{}}
}

func dectectRecursion(prog *Program) {
	for _, e := range prog.Exprs {
		fn, ok := e.(*FuncDeclExpr)
		if !ok {
			continue
		}
		if containsSelfCall(fn.Body, fn.Name.Lexeme) {
			fn.Recursion = true
		}
	}
}

func containsSelfCall(e Expr, fnName string) bool {
	switch n := e.(type) {
	case *CallExpr:
		if id, ok := n.Callee.(*Identifier); ok {
			if id.Name == fnName {
				return true
			}
		}
		for _, arg := range n.Args {
			if containsSelfCall(arg, fnName) {
				return true
			}
		}
		return false
	case *BlockExpr:
		for _, x := range n.Exprs {
			if containsSelfCall(x, fnName) {
				return true
			}
		}
	case *IfExpr:
		if containsSelfCall(n.Cond, fnName) ||
			containsSelfCall(n.Then, fnName) ||
			(n.Else != nil && containsSelfCall(n.Else, fnName)) {
			return true
		}
	case *MatchExpr:
		for _, arm := range n.Arms {
			if containsSelfCall(arm.Pattern, fnName) ||
				(arm.Guard != nil && containsSelfCall(arm.Guard, fnName)) ||
				containsSelfCall(arm.Body, fnName) {
				return true
			}
		}
	case *InfixExpr:
		return containsSelfCall(n.Left, fnName) ||
			containsSelfCall(n.Right, fnName)
	case *PrefixExpr:
		return containsSelfCall(n.Right, fnName)
	case *PipelineExpr:
		return containsSelfCall(n.Left, fnName) ||
			containsSelfCall(n.Right, fnName)
	case *TupleExpr:
		for _, t := range n.Elements {
			if containsSelfCall(t, fnName) {
				return true
			}
		}
	case *ListExpr:
		for _, el := range n.Elements {
			if containsSelfCall(el, fnName) {
				return true
			}
		}
	case *FieldAccessExpr:
		return containsSelfCall(n.Left, fnName)
	case *QualifiedExpr:
		return containsSelfCall(n.Left, fnName)
	case *VarDeclExpr:
		return containsSelfCall(n.Value, fnName)
	}
	return false
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
	if id, ok := left.(*Identifier); ok && p.cur().Kind == lexer.Equal {
		assignTok := p.eat()
		right := p.parseExpression(0)
		if right == nil {
			p.errorAt(assignTok, fmt.Sprintf("missing right-hand side for assignment to %s", id.Name))
		}
		return &AssignExpr{
			Name:  id,
			Value: right,
			Pos:   assignTok,
		}
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
			left = &FieldAccessExpr{
				Left:  left,
				Right: rightTok.Lexeme,
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
			p.errorAt(opTok, fmt.Sprintf("missing right-hand side after operator %q", opTok.Lexeme))
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
		for {
			if p.cur().Kind == lexer.Colon {
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
				continue
			}
			if p.cur().Kind == lexer.Dot {
				opTok := p.eat()
				fieldTok, ok := p.expect(lexer.Identifier)
				if !ok {
					return nil
				}
				expr = &FieldAccessExpr{
					Left:  expr,
					Right: fieldTok.Lexeme,
					Pos:   opTok,
				}
				continue
			}

			break
		}
		if p.cur().Kind == lexer.LeftParen {
			return p.parseCall(expr)
		}
		return expr
	case lexer.Int:
		p.eat()
		clean := lexer.StripNumericSeparators(tok.Lexeme)
		v, err := strconv.ParseInt(clean, 10, 64)
		if err != nil {
			p.errorAt(tok, fmt.Sprintf("invalid int literal %q", tok.Lexeme))
		}
		return &IntLiteral{Value: v, Raw: tok.Lexeme, Pos: tok}
	case lexer.Float:
		p.eat()
		clean := lexer.StripNumericSeparators(tok.Lexeme)
		f, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			p.errorAt(tok, fmt.Sprintf("invalid float literal %q", tok.Lexeme))
		}
		return &FloatLiteral{Value: f, Raw: tok.Lexeme, Pos: tok}
	case lexer.String:
		p.eat()
		value, err := strconv.Unquote(tok.Lexeme)
		if err != nil {
			p.errorAt(tok, "invalid string literal")
			return nil
		}
		return &StringLiteral{Value: value, Pos: tok}
	case lexer.Byte:
		p.eat()
		if len(tok.Lexeme) != 3 || tok.Lexeme[0] != '\'' || tok.Lexeme[2] != '\'' {
			p.errorAt(tok, fmt.Sprintf("invalid byte literal %q", tok.Lexeme))
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
			p.synchronize()
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
			p.errorAt(tok, fmt.Sprintf("missing expression after prefix %q", tok.Lexeme))
		}
		return &PrefixExpr{Operator: tok, Right: right}
	case lexer.KwPub:
		p.eat()
		switch p.cur().Kind {
		case lexer.KwFn:
			return p.parseFunc(true)
		case lexer.KwType:
			return p.recordTypeExpr(true)
		default:
			p.errorAt(p.cur(), "expected `fn` or `type` after `pub`")
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
	case lexer.LeftBracket:
		return p.parseList()
	case lexer.KwType:
		return p.recordTypeExpr(false)
	default:
		p.errorAt(tok, fmt.Sprintf("unexpected token %q", tok.Lexeme))
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
				p.errorAt(p.cur(), "invalid argument expression in call")
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
		p.synchronize()
		return nil
	}
	return &CallExpr{Callee: callee, Args: args, Pos: lparen}
}

func (p *Parser) parseVarDecl(mutable bool) Expr {
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
		p.synchronize()
		return nil
	}
	value := p.parseExpression(0)
	if value == nil {
		kind := "val"
		if mutable {
			kind = "mut"
		}
		p.errorAt(nameTok, fmt.Sprintf("missing initializer for %s %s", kind, nameTok.Lexeme))
	}
	return &VarDeclExpr{
		Mutable: mutable,
		Name:    nameTok,
		Type:    typeAnn,
		Value:   value,
	}
}

func (p *Parser) parseValDecl() Expr {
	return p.parseVarDecl(false)
}

func (p *Parser) parseMutDecl() Expr {
	return p.parseVarDecl(true)
}

func (p *Parser) parseFunc(pub bool) Expr {
	p.expect(lexer.KwFn)
	nameTok, ok := p.expect(lexer.Identifier)
	if !ok {
		return nil
	}

	if _, ok := p.expect(lexer.LeftParen); !ok {
		p.synchronize()
		return nil
	}

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

	if _, ok := p.expect(lexer.RightParen); !ok {
		p.synchronize()
		return nil
	}

	var retType Expr
	if p.cur().Kind != lexer.LeftBrace && p.cur().Kind != lexer.EndOfFile {
		retType = p.parseType()
	}

	var body Expr
	if p.cur().Kind == lexer.LeftBrace {
		body = p.parseBlock()
	}

	return &FuncDeclExpr{
		Pub:        pub,
		Name:       nameTok,
		Params:     params,
		Ret:        retType,
		Body:       body,
		Decorators: nil,
	}
}

func (p *Parser) parseBlock() Expr {
	if _, ok := p.expect(lexer.LeftBrace); !ok {
		p.synchronize()
		return nil
	}
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
	if _, ok := p.expect(lexer.RightBrace); !ok {
		p.synchronize()
	}
	return &BlockExpr{Exprs: exprs}
}

func (p *Parser) parseUse() Expr {
	start := p.eat()
	path := []string{}
	for {
		tok, ok := p.expect(lexer.Identifier)
		if !ok {
			p.errorAt(p.cur(), "expected identifier in use path")
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
				p.errorAt(p.cur(), "expected member in use {...}")
				return nil
			}
			members = append(members, memberTok.Lexeme)

			if p.cur().Kind == lexer.Comma {
				p.eat()
				continue
			}
			break
		}
		if _, ok := p.expect(lexer.RightBrace); !ok {
			p.errorAt(p.cur(), "expected '}' after members list")
			p.synchronize()
			return nil
		}
	}
	if p.cur().Kind == lexer.KwAs {
		p.eat()
		aliasTok, ok := p.expect(lexer.Identifier)
		if !ok {
			p.errorAt(p.cur(), "expected identifier after 'as'")
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
	start := p.cur()
	p.eat()
	cond := p.parseExpression(0)
	if cond == nil {
		p.errorAt(p.cur(), "expected condition after 'if'")
	}
	var thenExpr Expr
	var elseExpr Expr
	switch p.cur().Kind {
	case lexer.KwThen:
		p.eat()
		thenExpr = p.parseExpression(0)
		if thenExpr == nil {
			p.errorAt(p.cur(), "expected expression after then")
		}
		if p.cur().Kind != lexer.KwElse {
			p.errorAt(p.cur(), "expected 'else' after then-expression")
		}
		p.eat()
		if p.cur().Kind == lexer.LeftBrace {
			p.errorAt(p.cur(), "cannot use block-style else with expression-style then")
		}
		elseExpr = p.parseExpression(0)
		if elseExpr == nil {
			p.errorAt(p.cur(), "expected expression after else")
		}
	case lexer.LeftBrace:
		thenExpr = p.parseBlock()
		if p.cur().Kind == lexer.KwElse {
			p.eat()
			if p.cur().Kind != lexer.LeftBrace {
				p.errorAt(p.cur(), "block-style if requires block after else")
			}

			elseExpr = p.parseBlock()
		}
	default:
		p.errorAt(p.cur(), "expected 'then' or '{'")
		return nil
	}
	return &IfExpr{
		Pos:  start,
		Cond: cond,
		Then: thenExpr,
		Else: elseExpr,
	}
}

func (p *Parser) parseMatch() Expr {
	p.eat()
	value := p.parseExpression(0)
	if value == nil {
		p.errorAt(p.cur(), "expected expression after 'match'")
	}
	_, ok := p.expect(lexer.LeftBrace)
	if !ok {
		p.synchronize()
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
				p.errorAt(p.cur(), "expected pattern in match arm")
			}
		}
		var guard Expr
		if p.cur().Kind == lexer.KwIf {
			p.eat()
			guard = p.parseExpression(0)
			if guard == nil {
				p.errorAt(p.cur(), "expected guard expression after 'if'")
			}
		}
		_, ok := p.expect(lexer.RArrow)
		if !ok {
			return nil
		}
		body := p.parseExpression(0)
		if body == nil {
			p.errorAt(p.cur(), "expected body expression in match arm")
		}
		arms = append(arms, &MatchArm{
			Pattern: pattern,
			Guard:   guard,
			Body:    body,
		})
	}
	_, ok = p.expect(lexer.RightBrace)
	if !ok {
		p.synchronize()
		return nil
	}
	return &MatchExpr{
		Value: value,
		Arms:  arms,
	}
}

func (p *Parser) parseList() Expr {
	start := p.eat()
	elements := []Expr{}
	for p.cur().Kind != lexer.RightBracket && p.cur().Kind != lexer.EndOfFile {
		elem := p.parseExpression(0)
		if elem == nil {
			p.errorAt(p.cur(), "invalid list element")
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
	case lexer.Identifier:
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
		if _, ok := p.expect(lexer.RightParen); !ok {
			p.synchronize()
			return nil
		}
		return &TupleTypeExpr{Types: types, Pos: tok}
	default:
		p.errorAt(tok, fmt.Sprintf("expected type, got %q (%v)", tok.Lexeme, tok.Kind))
		return nil
	}
}

func (p *Parser) recordTypeExpr(pub bool) Expr {
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
				p.errorAt(fieldTok, fmt.Sprintf("expected ':' after field name %s", fieldTok.Lexeme))
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
		if _, ok := p.expect(lexer.RightBrace); !ok {
			p.synchronize()
			return nil
		}
		body = &RecordTypeExpr{Name: nameTok, Fields: fields, Pos: nameTok}
	}
	return &TypeDeclExpr{
		Pub:  pub,
		Name: nameTok,
		Body: body,
		Pos:  nameTok,
	}
}

func (p *Parser) parseDecorators() []Decorator {
	decorators := []Decorator{}
	for p.cur().Kind == lexer.At {
		p.eat()
		nameTok, ok := p.expect(lexer.Identifier)
		if !ok {
			p.errorAt(p.cur(), "expected decorator name after '@'")
			break
		}
		args := []Expr{}
		if p.cur().Kind == lexer.LeftParen {
			p.eat()
			for p.cur().Kind != lexer.RightParen && p.cur().Kind != lexer.EndOfFile {
				arg := p.parseExpression(0)
				if arg == nil {
					p.errorAt(p.cur(), "invalid decorator argument")
				} else {
					args = append(args, arg)
				}
				if p.cur().Kind == lexer.Comma {
					p.eat()
				}
			}
			p.expect(lexer.RightParen)
		}
		decorators = append(decorators, Decorator{Name: nameTok.Lexeme, Args: args, Pos: nameTok})
	}
	return decorators
}
