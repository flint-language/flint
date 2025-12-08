package typechecker

import (
	"fmt"
	"strings"

	"flint/internal/lexer"
	"flint/internal/parser"
)

type TypeChecker struct {
	errors []string
	env    *Env
	ctx    Context
}

func New() *TypeChecker {
	return &TypeChecker{
		errors: []string{},
		env:    NewEnv(nil),
		ctx:    TopLevel,
	}
}

func (tc *TypeChecker) CheckExpr(expr parser.Expr) (*Type, error) {
	ty := tc.Check(expr)
	if len(tc.errors) > 0 {
		err := fmt.Errorf("%s", tc.errors[0])
		tc.errors = tc.errors[1:]
		return &Type{TKind: TyError}, err
	}
	return ty, nil
}

func (tc *TypeChecker) Check(expr parser.Expr) *Type {
	if tc.ctx == TopLevel {
		switch expr.(type) {
		case *parser.VarDeclExpr, *parser.IfExpr,
			*parser.MatchExpr, *parser.PipelineExpr:
			tc.error(lexer.Token{}, fmt.Sprintf("%T is not allowed at top-level; must be inside a function/block", expr))
			return &Type{TKind: TyError}
		}
	}
	switch e := expr.(type) {
	case *parser.IntLiteral:
		return &Type{TKind: TyInt}
	case *parser.FloatLiteral:
		return &Type{TKind: TyFloat}
	case *parser.BoolLiteral:
		return &Type{TKind: TyBool}
	case *parser.StringLiteral:
		return &Type{TKind: TyString}
	case *parser.ByteLiteral:
		return &Type{TKind: TyByte}
	case *parser.PrefixExpr:
		return tc.visitPrefix(e)
	case *parser.InfixExpr:
		return tc.visitInfix(e)
	case *parser.Identifier:
		return tc.visitIdentifier(e)
	case *parser.VarDeclExpr:
		return tc.visitVarDecl(e)
	case *parser.FuncDeclExpr:
		oldCtx := tc.ctx
		tc.ctx = FunctionBody
		ty := tc.visitFuncDecl(e)
		tc.ctx = oldCtx
		return ty
	case *parser.CallExpr:
		oldCtx := tc.ctx
		tc.ctx = FunctionBody
		ty := tc.visitCall(e)
		tc.ctx = oldCtx
		return ty
	case *parser.BlockExpr:
		return tc.visitBlock(e)
	case *parser.UseExpr:
		return tc.visitUse(e)
	case *parser.QualifiedExpr:
		return tc.visitQualified(e)
	case *parser.IfExpr:
		return tc.visitIf(e)
	case *parser.MatchExpr:
		return tc.visitMatch(e)
	case *parser.PipelineExpr:
		return tc.visitPipeline(e)
	case *parser.ListExpr:
		return tc.visitList(e, nil)
	case *parser.AssignExpr:
		return tc.visitAssign(e)
	default:
		return &Type{TKind: TyError}
	}
}

func (tc *TypeChecker) visitIdentifier(id *parser.Identifier) *Type {
	ty, ok := tc.env.Get(id.Name)
	if !ok {
		tc.error(id.Pos, "undefined variable: "+id.Name)
		return &Type{TKind: TyError}
	}
	return ty
}

func (tc *TypeChecker) visitVarDecl(d *parser.VarDeclExpr) *Type {
	var varTy *Type
	if d.Value != nil {
		switch expr := d.Value.(type) {
		case *parser.ListExpr:
			varTy = tc.visitList(expr, nil)
		default:
			varTy = tc.Check(expr)
		}
		if varTy == nil || varTy.TKind == TyError {
			tc.error(d.Name, fmt.Sprintf("cannot infer type for %s '%s'",
				func() string {
					if d.Mutable {
						return "mut"
					}
					return "val"
				}(), d.Name.Lexeme))
			return &Type{TKind: TyError}
		}
	}
	if d.Type != nil {
		declTy := tc.resolveType(d.Type)
		if varTy != nil && !declTy.Equal(varTy) {
			tc.error(d.Name, fmt.Sprintf(
				"type mismatch in %s '%s': expected %s, got %s",
				func() string {
					if d.Mutable {
						return "mut"
					}
					return "val"
				}(), d.Name.Lexeme, declTy.String(), varTy.String()))
			tc.env.Set(d.Name.Lexeme, declTy)
			return &Type{TKind: TyError}
		}
		tc.env.SetVar(d.Name.Lexeme, declTy, d.Mutable)
		return declTy
	}
	tc.env.SetVar(d.Name.Lexeme, varTy, d.Mutable)
	return varTy
}

func (tc *TypeChecker) visitFuncDecl(fn *parser.FuncDeclExpr) *Type {
	paramTypes := make([]*Type, len(fn.Params))
	for i, p := range fn.Params {
		if p.Type == nil {
			tc.error(p.Name,
				fmt.Sprintf("parameter '%s' missing type annotation", p.Name.Lexeme))
			return &Type{TKind: TyError}
		}

		pt := tc.resolveType(p.Type)
		if pt.TKind == TyError {
			return &Type{TKind: TyError}
		}
		paramTypes[i] = pt
	}
	var retType *Type
	if fn.Ret != nil {
		retType = tc.resolveType(fn.Ret)
		if retType.TKind == TyError {
			return &Type{TKind: TyError}
		}
	} else {
		retType = &Type{TKind: TyNil}
	}
	fnType := &Type{
		TKind:  TyFunc,
		Params: paramTypes,
		Ret:    retType,
	}
	tc.env.Set(fn.Name.Lexeme, fnType)
	oldEnv := tc.env
	tc.env = NewEnv(oldEnv)
	for i, p := range fn.Params {
		tc.env.Set(p.Name.Lexeme, paramTypes[i])
	}
	if fn.Body != nil {
		bodyTy := tc.Check(fn.Body)
		if fn.Ret != nil && !bodyTy.Equal(retType) {
			tc.error(fn.Name,
				fmt.Sprintf("function %s annotated return %s but body has type %s",
					fn.Name.Lexeme, retType.String(), bodyTy.String()))
			tc.env.Set(fn.Name.Lexeme, &Type{TKind: TyError})
			tc.env = oldEnv
			return &Type{TKind: TyError}
		}
		if fn.Ret == nil {
			fnType.Ret = bodyTy
		}
	}
	tc.env = oldEnv
	return fnType
}

func (tc *TypeChecker) visitCall(c *parser.CallExpr) *Type {
	calleeTy := tc.Check(c.Callee)
	if calleeTy.TKind != TyFunc {
		tc.error(c.Pos, fmt.Sprintf("attempt to call non-function value of type %s", calleeTy.String()))
		return &Type{TKind: TyError}
	}
	if len(c.Args) != len(calleeTy.Params) {
		tc.error(c.Pos, fmt.Sprintf("wrong number of arguments: expected %d, got %d",
			len(calleeTy.Params), len(c.Args)))
		return &Type{TKind: TyError}
	}
	for i, a := range c.Args {
		argTy := tc.Check(a)
		if !argTy.Equal(calleeTy.Params[i]) {
			tc.error(c.Pos, fmt.Sprintf("argument %d expected %s, got %s",
				i, calleeTy.Params[i].String(), argTy.String()))
			return &Type{TKind: TyError}
		}
	}
	return calleeTy.Ret
}

func (tc *TypeChecker) visitBlock(b *parser.BlockExpr) *Type {
	old := tc.env
	tc.env = NewEnv(old)
	last := &Type{TKind: TyNil}
	for _, ex := range b.Exprs {
		last = tc.Check(ex)
	}
	tc.env = old
	return last
}

func (tc *TypeChecker) visitPrefix(e *parser.PrefixExpr) *Type {
	arg := tc.Check(e.Right)
	sig, ok := unaryOps[e.Operator.Kind]
	if !ok {
		tc.error(e.Operator, "unknown unary operator")
		return &Type{TKind: TyError}
	}
	if arg.TKind != sig.Arg.TKind {
		tc.error(e.Operator,
			fmt.Sprintf("invalid operand type for '%s': %s",
				e.Operator.Lexeme, arg.String()))
		return &Type{TKind: TyError}
	}
	out := sig.Out
	return &out
}

func (tc *TypeChecker) visitInfix(e *parser.InfixExpr) *Type {
	left := tc.Check(e.Left)
	right := tc.Check(e.Right)
	sigs, ok := binOps[e.Operator.Kind]
	if !ok {
		tc.error(e.Operator, "unknown operator")
		return &Type{TKind: TyError}
	}
	for _, sig := range sigs {
		if left.TKind == sig.Left.TKind && right.TKind == sig.Right.TKind {
			out := sig.Out
			return &out
		}
	}
	tc.error(e.Operator,
		fmt.Sprintf("invalid operands for '%s': %s and %s",
			e.Operator.Lexeme, left.String(), right.String()))

	return &Type{TKind: TyError}
}

func (tc *TypeChecker) visitUse(u *parser.UseExpr) *Type {
	modEnv, ok := getModule(u.Path)
	if !ok {
		tc.error(u.Pos, fmt.Sprintf("cannot find module %s", strings.Join(u.Path, "/")))
		return &Type{TKind: TyError}
	}
	if len(u.Members) == 0 {
		name := u.Alias
		if name == "" && len(u.Path) > 0 {
			name = u.Path[len(u.Path)-1]
		}
		tc.env.modules[name] = modEnv
	} else {
		for _, m := range u.Members {
			ty, ok := modEnv.Get(m)
			if !ok {
				tc.error(u.Pos, fmt.Sprintf("module %s has no member %s", strings.Join(u.Path, "/"), m))
				continue
			}
			tc.env.Set(m, ty)
		}
	}
	return &Type{TKind: TyNil}
}

func (tc *TypeChecker) visitQualified(q *parser.QualifiedExpr) *Type {
	leftIdent, ok := q.Left.(*parser.Identifier)
	if !ok {
		tc.error(q.Pos, "expected module identifier on the left of ':'")
		return &Type{TKind: TyError}
	}
	modEnv, ok := tc.env.modules[leftIdent.Name]
	if !ok {
		tc.error(q.Pos, fmt.Sprintf("unknown module: %s", leftIdent.Name))
		return &Type{TKind: TyError}
	}
	ty, ok := modEnv.Get(q.Right.Lexeme)
	if !ok {
		tc.error(q.Pos, fmt.Sprintf("module %s has no member %s", leftIdent.Name, q.Right.Lexeme))
		return &Type{TKind: TyError}
	}
	return ty
}

func (tc *TypeChecker) visitIf(i *parser.IfExpr) *Type {
	condTy := tc.Check(i.Cond)
	if condTy.TKind != TyBool {
		tc.error(i.Pos, fmt.Sprintf(
			"if condition must be Bool, got %s", condTy.String()))
		return &Type{TKind: TyError}
	}
	thenTy := tc.Check(i.Then)
	if i.Else != nil {
		elseTy := tc.Check(i.Else)
		if !thenTy.Equal(elseTy) {
			tc.error(i.Pos, fmt.Sprintf(
				"then branch has type %s but else branch has type %s",
				thenTy.String(), elseTy.String()))
			return &Type{TKind: TyError}
		}
	}
	return thenTy
}

func (tc *TypeChecker) visitMatch(m *parser.MatchExpr) *Type {
	valueTy := tc.Check(m.Value)
	if valueTy.TKind == TyError {
		return &Type{TKind: TyError}
	}
	var armType *Type
	for _, arm := range m.Arms {
		oldEnv := tc.env
		tc.env = NewEnv(oldEnv)
		switch p := arm.Pattern.(type) {
		case *parser.Identifier:
			tc.env.Set(p.Name, valueTy)
			patternTy := valueTy
			_ = patternTy
		default:
			patternTy := tc.Check(arm.Pattern)
			if !patternTy.Equal(valueTy) {
				tc.error(arm.Pos, fmt.Sprintf(
					"pattern type %s does not match value type %s",
					patternTy.String(), valueTy.String()))
				tc.env = oldEnv
				return &Type{TKind: TyError}
			}
		}
		if arm.Guard != nil {
			guardTy := tc.Check(arm.Guard)
			if guardTy.TKind != TyBool {
				tc.error(arm.Pos, fmt.Sprintf(
					"guard must be Bool, got %s", guardTy.String()))
				tc.env = oldEnv
				return &Type{TKind: TyError}
			}
		}
		bodyTy := tc.Check(arm.Body)
		if armType == nil {
			armType = bodyTy
		} else if !armType.Equal(bodyTy) {
			tc.error(arm.Pos, fmt.Sprintf(
				"match arm has type %s, expected %s",
				bodyTy.String(), armType.String()))
			tc.env = oldEnv
			return &Type{TKind: TyError}
		}
		tc.env = oldEnv
	}
	if armType == nil {
		return &Type{TKind: TyNil}
	}
	return armType
}

func (tc *TypeChecker) visitPipeline(p *parser.PipelineExpr) *Type {
	leftTy := tc.Check(p.Left)
	if leftTy.TKind == TyError {
		return &Type{TKind: TyError}
	}
	switch r := p.Right.(type) {
	case *parser.Identifier:
		fnTy, ok := tc.env.Get(r.Name)
		if !ok {
			tc.error(r.Pos, "undefined function: "+r.Name)
			return &Type{TKind: TyError}
		}
		if fnTy.TKind != TyFunc || len(fnTy.Params) == 0 {
			tc.error(r.Pos, fmt.Sprintf("cannot pipe to non-function or function with no parameters: %s", r.Name))
			return &Type{TKind: TyError}
		}
		if !fnTy.Params[0].Equal(leftTy) {
			tc.error(r.Pos, fmt.Sprintf("type mismatch in pipeline: expected %s, got %s",
				fnTy.Params[0].String(), leftTy.String()))
			return &Type{TKind: TyError}
		}
		return fnTy.Ret
	case *parser.CallExpr:
		args := append([]parser.Expr{p.Left}, r.Args...)
		call := &parser.CallExpr{
			Callee: r.Callee,
			Args:   args,
			Pos:    r.Pos,
		}
		return tc.Check(call)
	default:
		tc.error(p.Pos, "right side of pipeline must be a function or call")
		return &Type{TKind: TyError}
	}
}

func (tc *TypeChecker) visitList(l *parser.ListExpr, annotated *Type) *Type {
	if len(l.Elements) == 0 {
		if annotated != nil {
			return annotated
		}
		return &Type{TKind: TyList, Elem: &Type{TKind: TyNil}}
	}
	var expected *Type
	if annotated != nil && annotated.TKind == TyList {
		expected = annotated.Elem
	} else {
		first := l.Elements[0]
		switch tup := first.(type) {
		case *parser.TupleExpr:
			tElems := make([]*Type, len(tup.Elements))
			for i, e := range tup.Elements {
				subTy := tc.Check(e)
				if subTy == nil || subTy.TKind == TyError {
					tc.error(l.Pos, "cannot infer element type for tuple in list")
					return &Type{TKind: TyError}
				}
				tElems[i] = subTy
			}
			expected = &Type{TKind: TyTuple, TElems: tElems}
		default:
			expected = tc.Check(first)
			if expected == nil || expected.TKind == TyError {
				tc.error(l.Pos, "cannot infer element type for list (first element error)")
				return &Type{TKind: TyError}
			}
		}
	}
	if expected.TKind == TyTuple {
		for i, e := range l.Elements {
			tup, ok := e.(*parser.TupleExpr)
			if !ok {
				tc.error(l.Pos, fmt.Sprintf(
					"element %d: expected tuple %s, got %s",
					i+1, expected.String(), tc.Check(e).String()))
				return &Type{TKind: TyError}
			}
			if len(tup.Elements) != len(expected.TElems) {
				tc.error(tup.Pos, fmt.Sprintf(
					"element %d: expected tuple of length %d, got %d",
					i+1, len(expected.TElems), len(tup.Elements)))
				return &Type{TKind: TyError}
			}
			for k, sub := range tup.Elements {
				subTy := tc.Check(sub)
				if !expected.TElems[k].Equal(subTy) {
					tc.error(l.Pos, fmt.Sprintf(
						"element %d.%d: expected %s, got %s",
						i+1, k+1, expected.TElems[k].String(), subTy.String()))
					return &Type{TKind: TyError}
				}
			}
		}
	} else {
		for i, e := range l.Elements {
			ty := tc.Check(e)
			if !expected.Equal(ty) {
				tc.error(l.Pos, fmt.Sprintf(
					"element %d type %s does not match expected type %s",
					i+1, ty.String(), expected.String()))
				return &Type{TKind: TyError}
			}
		}
	}
	return &Type{TKind: TyList, Elem: expected}
}

func (tc *TypeChecker) visitAssign(a *parser.AssignExpr) *Type {
	varInfo, ok := tc.env.GetVar(a.Name.Name)
	if !ok {
		tc.error(a.Pos, fmt.Sprintf("undefined variable '%s'", a.Name.Name))
		return &Type{TKind: TyError}
	}
	if !varInfo.Mutable {
		tc.error(a.Pos, fmt.Sprintf("cannot assign to immutable variable '%s'", a.Name.Name))
		return &Type{TKind: TyError}
	}
	valueTy := tc.Check(a.Value)
	if !varInfo.Ty.Equal(valueTy) {
		tc.error(a.Pos, fmt.Sprintf(
			"type mismatch in assignment to '%s': expected %s, got %s",
			a.Name.Name, varInfo.Ty.String(), valueTy.String()))
		return &Type{TKind: TyError}
	}
	tc.env.SetVar(a.Name.Name, valueTy, true)
	return valueTy
}

// TODO: Add records
