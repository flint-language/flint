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

func TestUnsignedCoercion(t *testing.T) {
	ty, err := typeOf(t, `
fun test() U32 { 
	val x: U32 = 10 + 1
	x 
}
`)
	if err != nil {
		t.Fatal(err)
	}
	if ty.Ret.TKind != TyU32 {
		t.Fatalf("expected U32, got %s", ty)
	}
}

func TestFloatCoercion(t *testing.T) {
	ty, err := typeOf(t, `
fun test() F64 { 
	val x: F64 = 1.0 +. 2.0 
	x
}
`)
	if err != nil {
		t.Fatal(err)
	}
	if ty.Ret.TKind != TyF64 {
		t.Fatalf("expected F64, got %s", ty)
	}
}

func TestNestedInfixCoercion(t *testing.T) {
	ty, err := typeOf(t, `fun test() U32 { val x: U32 = 1 + 2 + 3 }`)
	if err != nil {
		t.Fatal(err)
	}
	if ty.Ret.TKind != TyU32 {
		t.Fatalf("expected U32, got %s", ty)
	}
}

func TestInvalidAdditionMixTypes(t *testing.T) {
	_, err := typeOf(t, `val x: Int = 1 + 2.5`)
	if err == nil {
		t.Fatal("expected type error for mixing Int + Float")
	}
}

func TestUnsignedAddition(t *testing.T) {
	ty, err := typeOf(t, `fun test() U64 { val x: U64 = 10 + 20 }`)
	if err != nil {
		t.Fatal(err)
	}
	if ty.Ret.TKind != TyU64 {
		t.Fatalf("expected U64, got %s", ty)
	}
}

func TestFloatArithmetic(t *testing.T) {
	ty, err := typeOf(t, `fun test() F32 { val x: F32 = 1.5 +. 2.5 }`)
	if err != nil {
		t.Fatal(err)
	}
	if ty.Ret.TKind != TyF32 {
		t.Fatalf("expected F32, got %s", ty)
	}
}
