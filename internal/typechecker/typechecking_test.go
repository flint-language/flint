package typechecker

import (
	"flint/internal/lexer"
	"flint/internal/parser"
	"testing"
)

func typeOf(t *testing.T, src string) (*Type, error) {
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

	prog, errs := parser.ParseProgram(tokens)
	if len(errs) > 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	if len(prog.Exprs) == 0 {
		t.Fatalf("no expressions parsed")
	}

	tc := New()
	return tc.CheckExpr(prog.Exprs[0])
}

func TestTypeIntLiteral(t *testing.T) {
	ty, err := typeOf(t, "123")
	if err != nil {
		t.Fatal(err)
	}
	if ty.TKind != TyInt {
		t.Fatalf("expected Int, got %s", ty)
	}
}

func TestTypeStringLiteral(t *testing.T) {
	ty, err := typeOf(t, `"hello"`)
	if err != nil {
		t.Fatal(err)
	}
	if ty.TKind != TyString {
		t.Fatalf("expected String, got %s", ty)
	}
}

func TestTypeBoolLiteral(t *testing.T) {
	ty, err := typeOf(t, "True")
	if err != nil {
		t.Fatal(err)
	}
	if ty.TKind != TyBool {
		t.Fatalf("expected Bool, got %s", ty)
	}
}

func TestTypeAddition(t *testing.T) {
	ty, err := typeOf(t, "1 + 2")
	if err != nil {
		t.Fatal(err)
	}
	if ty.TKind != TyInt {
		t.Fatalf("expected Int, got %s", ty)
	}
}

func TestTypeMismatch(t *testing.T) {
	_, err := typeOf(t, `1 + "hello"`)
	if err == nil {
		t.Fatal("expected type error, got none")
	}
}

func TestFunctionDeclaration(t *testing.T) {
	ty, err := typeOf(t, `
fun add(x: Int, y: Int) Int {
	x + y
}
`)
	if err != nil {
		t.Fatal(err)
	}

	if ty.TKind != TyFunc {
		t.Fatalf("expected function type, got %s", ty)
	}
}

func TestFunctionBadReturn(t *testing.T) {
	_, err := typeOf(t, `
fun bad(x: Int) Bool {
	x
}
`)
	if err == nil {
		t.Fatal("expected type error for return mismatch")
	}
}
