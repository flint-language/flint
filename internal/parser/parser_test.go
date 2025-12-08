package parser

import (
	"flint/internal/lexer"
	"testing"
)

func parseSrc(t *testing.T, src string) (*Program, []string) {
	t.Helper()

	l := lexer.New(src, "test.flint")
	tokens := []lexer.Token{}
	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Kind == lexer.EndOfFile {
			break
		}
	}

	prog, errs := ParseProgram(tokens)
	return prog, errs
}

func TestParseIntLiteral(t *testing.T) {
	prog, errs := parseSrc(t, "123")

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	if len(prog.Exprs) != 1 {
		t.Fatalf("expected 1 expr, got %d", len(prog.Exprs))
	}

	n, ok := prog.Exprs[0].(*IntLiteral)
	if !ok {
		t.Fatalf("expected IntLiteral, got %T", prog.Exprs[0])
	}

	if n.Value != 123 {
		t.Fatalf("expected 123, got %d", n.Value)
	}
}

func TestParseStringLiteral(t *testing.T) {
	prog, errs := parseSrc(t, `"hello"`)

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	s, ok := prog.Exprs[0].(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral, got %T", prog.Exprs[0])
	}

	if s.Value != "hello" {
		t.Fatalf("expected hello, got %q", s.Value)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	prog, errs := parseSrc(t, "1 + 2 * 3")

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	root, ok := prog.Exprs[0].(*InfixExpr)
	if !ok {
		t.Fatalf("expected InfixExpr, got %T", prog.Exprs[0])
	}

	if root.Operator.Lexeme != "+" {
		t.Fatalf("expected '+', got %s", root.Operator.Lexeme)
	}

	right := root.Right.(*InfixExpr)
	if right.Operator.Lexeme != "*" {
		t.Fatalf("expected '*', got %s", right.Operator.Lexeme)
	}
}

func TestFunctionCall(t *testing.T) {
	prog, errs := parseSrc(t, "add(1, 2)")

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	call, ok := prog.Exprs[0].(*CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", prog.Exprs[0])
	}

	if len(call.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(call.Args))
	}
}

func TestValDecl(t *testing.T) {
	prog, errs := parseSrc(t, "val x = 10")

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	v, ok := prog.Exprs[0].(*VarDeclExpr)
	if !ok {
		t.Fatalf("expected ValDeclExpr, got %T", prog.Exprs[0])
	}

	if v.Name.Lexeme != "x" {
		t.Fatalf("expected name 'x', got %s", v.Name.Lexeme)
	}
}

func TestFunctionDecl(t *testing.T) {
	src := `
fn add(x: Int, y: Int) Int {
	x + y
}
`
	prog, errs := parseSrc(t, src)

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	fn, ok := prog.Exprs[0].(*FuncDeclExpr)
	if !ok {
		t.Fatalf("expected FuncDeclExpr, got %T", prog.Exprs[0])
	}

	if fn.Name.Lexeme != "add" {
		t.Fatalf("wrong function name %s", fn.Name.Lexeme)
	}

	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
}

func TestIfExpression(t *testing.T) {
	src := `if x then 1 else 2`

	prog, errs := parseSrc(t, src)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	ifx, ok := prog.Exprs[0].(*IfExpr)
	if !ok {
		t.Fatalf("expected IfExpr, got %T", prog.Exprs[0])
	}

	if ifx.Then == nil || ifx.Else == nil {
		t.Fatalf("missing then or else branch")
	}
}

func TestMissingRHS(t *testing.T) {
	_, errs := parseSrc(t, `1 +`)

	if len(errs) == 0 {
		t.Fatal("expected error but got none")
	}
}

func TestBadFunctionSyntax(t *testing.T) {
	_, errs := parseSrc(t, `fn (x) { x }`)

	if len(errs) == 0 {
		t.Fatal("expected error for missing function name")
	}
}
