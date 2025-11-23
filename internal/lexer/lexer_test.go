package lexer_test

import (
	lx "flint/internal/lexer"
	"testing"
)

func TestLexerBasicToken(t *testing.T) {
	input := `mut x = 10
		y = x + 2
		// comment
		"hi"
		'a'
		3.14
		1..5`

	tests := []struct {
		Kind   lx.TokenKind
		Lexeme string
	}{
		{lx.KwMut, "mut"},
		{lx.Identifier, "x"},
		{lx.Equal, "="},
		{lx.Int, "10"},
		{lx.Identifier, "y"},
		{lx.Equal, "="},
		{lx.Identifier, "x"},
		{lx.Plus, "+"},
		{lx.Int, "2"},
		{lx.Comment, "// comment"},
		{lx.String, `"hi"`},
		{lx.Byte, `'a'`},
		{lx.Float, "3.14"},
		{lx.Int, "1"},
		{lx.DotDot, ".."},
		{lx.Int, "5"},
		{lx.EndOfFile, ""},
	}

	lexer := lx.New(input, "test.flint")
	for i, tt := range tests {
		tok := lexer.Next()

		if tok.Kind != tt.Kind {
			t.Fatalf("test %d: expected kind %v, got %v (lexeme=%q)",
				i, tt.Kind, tok.Kind, tok.Lexeme)
		}

		if tok.Lexeme != tt.Lexeme {
			t.Fatalf("test %d: expected lexeme %q, got %q",
				i, tt.Lexeme, tok.Lexeme)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	input := "foo bar _x $y z' tail'"
	expected := []string{"foo", "bar", "_x", "$y", "z'", "tail'"}

	lexer := lx.New(input, "test.flint")

	for _, want := range expected {
		tok := lexer.Next()
		if tok.Kind != lx.Identifier {
			t.Fatalf("expected Identifier, got %v", tok.Kind)
		}
		if tok.Lexeme != want {
			t.Fatalf("expected %q, got %q", want, tok.Lexeme)
		}
	}
}

func TestNumbers(t *testing.T) {
	input := "123 4_567 3.1415 10.0"
	kinds := []lx.TokenKind{lx.Int, lx.Int, lx.Float, lx.Float}
	lexemes := []string{"123", "4_567", "3.1415", "10.0"}

	lexer := lx.New(input, "numbers.flint")

	for i := range lexemes {
		tok := lexer.Next()
		if tok.Kind != kinds[i] {
			t.Fatalf("expected %v, got %v", kinds[i], tok.Kind)
		}
		if tok.Lexeme != lexemes[i] {
			t.Fatalf("expected %q, got %q", lexemes[i], tok.Lexeme)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	lexer := lx.New(`"hello\nworld"`, "str.flint")
	tok := lexer.Next()

	if tok.Kind != lx.String {
		t.Fatalf("expected String, got %v", tok.Kind)
	}
	if tok.Lexeme != `"hello\nworld"` {
		t.Fatalf("unexpected lexeme: %q", tok.Lexeme)
	}
}

func TestByteLiteral(t *testing.T) {
	lexer := lx.New(`'\n'`, "byte.flint")
	tok := lexer.Next()

	if tok.Kind != lx.Byte {
		t.Fatalf("expected Byte, got %v", tok.Kind)
	}
	if tok.Lexeme != `'\n'` {
		t.Fatalf("unexpected lexeme: %q", tok.Lexeme)
	}
}
