package lexer

import "testing"

func TestLexerBasicToken(t *testing.T) {
	input := `mut x = 10
y = x + 2
// comment
"hi"
'a'
3.14
1..5`

	tests := []struct {
		Kind   TokenKind
		Lexeme string
	}{
		{KwMut, "mut"},
		{Identifier, "x"},
		{Equal, "="},
		{Int, "10"},
		{Identifier, "y"},
		{Equal, "="},
		{Identifier, "x"},
		{Plus, "+"},
		{Int, "2"},
		{Comment, "// comment"},
		{String, `"hi"`},
		{Byte, `'a'`},
		{Float, "3.14"},
		{Int, "1"},
		{DotDot, ".."},
		{Int, "5"},
		{EndOfFile, ""},
	}

	lexer := New(input, "test.flint")

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
	input := "foo bar _x $y z' tail"
	expected := []string{"foo", "bar", "_x", "$y", "z'", "tail"}

	lexer := New(input, "identifiers.flint")

	for i, want := range expected {
		tok := lexer.Next()

		if tok.Kind != Identifier {
			t.Fatalf("test %d: expected Identifier, got %v", i, tok.Kind)
		}

		if tok.Lexeme != want {
			t.Fatalf("test %d: expected %q, got %q", i, want, tok.Lexeme)
		}
	}

	if tok := lexer.Next(); tok.Kind != EndOfFile {
		t.Fatalf("expected EOF, got %v", tok.Kind)
	}
}

func TestNumbers(t *testing.T) {
	input := "123 4_567 3.14159 10.0"
	kinds := []TokenKind{Int, Int, Float, Float}
	lexemes := []string{"123", "4_567", "3.14159", "10.0"}

	lexer := New(input, "numbers.flint")

	for i := range lexemes {
		tok := lexer.Next()

		if tok.Kind != kinds[i] {
			t.Fatalf("test %d: expected %v, got %v", i, kinds[i], tok.Kind)
		}

		if tok.Lexeme != lexemes[i] {
			t.Fatalf("test %d: expected %q, got %q", i, lexemes[i], tok.Lexeme)
		}
	}

	if tok := lexer.Next(); tok.Kind != EndOfFile {
		t.Fatalf("expected EOF, got %v", tok.Kind)
	}
}

func TestStringLiteral(t *testing.T) {
	lexer := New(`"hello\nworld"`, "strings.flint")
	tok := lexer.Next()

	if tok.Kind != String {
		t.Fatalf("expected String, got %v", tok.Kind)
	}

	if tok.Lexeme != `"hello\nworld"` {
		t.Fatalf("unexpected lexeme: %q", tok.Lexeme)
	}

	if tok := lexer.Next(); tok.Kind != EndOfFile {
		t.Fatalf("expected EOF, got %v", tok.Kind)
	}
}

func TestByteLiteral(t *testing.T) {
	lexer := New(`'\n'`, "byte.flint")
	tok := lexer.Next()

	if tok.Kind != Byte {
		t.Fatalf("expected Byte, got %v", tok.Kind)
	}

	if tok.Lexeme != `'\n'` {
		t.Fatalf("unexpected lexeme: %q", tok.Lexeme)
	}

	if tok := lexer.Next(); tok.Kind != EndOfFile {
		t.Fatalf("expected EOF, got %v", tok.Kind)
	}
}

func TestKeywordsVsIdentifiers(t *testing.T) {
	input := "mut mutable fn function val letter"
	lexer := New(input, "keywords.flint")

	tests := []struct {
		kind   TokenKind
		lexeme string
	}{
		{KwMut, "mut"},
		{Identifier, "mutable"},
		{KwFn, "fn"},
		{Identifier, "function"},
		{KwVal, "val"},
		{Identifier, "letter"},
		{EndOfFile, ""},
	}

	for i, tt := range tests {
		tok := lexer.Next()

		if tok.Kind != tt.kind {
			t.Fatalf("test %d: expected %v, got %v", i, tt.kind, tok.Kind)
		}

		if tok.Lexeme != tt.lexeme {
			t.Fatalf("test %d: expected %q, got %q", i, tt.lexeme, tok.Lexeme)
		}
	}
}

func TestOperators(t *testing.T) {
	input := "+ - * / +. -. *. /. == != <= <=. >=. >= < > && || ! <>"
	lexer := New(input, "operators.flint")

	tests := []struct {
		kind TokenKind
	}{
		{Plus},
		{Minus},
		{Star},
		{Slash},
		{PlusDot},
		{MinusDot},
		{StarDot},
		{SlashDot},
		{EqualEqual},
		{NotEqual},
		{LessEqual},
		{LessEqualDot},
		{GreaterEqualDot},
		{GreaterEqual},
		{Less},
		{Greater},
		{AmperAmper},
		{VbarVbar},
		{Bang},
		{LtGt},
		{EndOfFile},
	}

	for i, tt := range tests {
		tok := lexer.Next()

		if tok.Kind != tt.kind {
			t.Fatalf("test %d: expected %v, got %v (%q)",
				i, tt.kind, tok.Kind, tok.Lexeme)
		}
	}
}

func TestRanges(t *testing.T) {
	input := "1..5 3.14 10."
	lexer := New(input, "ranges.flint")

	tests := []struct {
		kind   TokenKind
		lexeme string
	}{
		{Int, "1"},
		{DotDot, ".."},
		{Int, "5"},
		{Float, "3.14"},
		{Float, "10."},
		{EndOfFile, ""},
	}

	for i, tt := range tests {
		tok := lexer.Next()

		if tok.Kind != tt.kind {
			t.Fatalf("test %d: expected %v, got %v", i, tt.kind, tok.Kind)
		}

		if tok.Lexeme != tt.lexeme {
			t.Fatalf("test %d: expected %q, got %q", i, tt.lexeme, tok.Lexeme)
		}
	}
}

func TestCommentEdgeCases(t *testing.T) {
	input := "//first\n//second\nx"
	lexer := New(input, "comments.flint")

	tests := []struct {
		kind   TokenKind
		lexeme string
	}{
		{Comment, "//first"},
		{Comment, "//second"},
		{Identifier, "x"},
		{EndOfFile, ""},
	}

	for i, tt := range tests {
		tok := lexer.Next()

		if tok.Kind != tt.kind {
			t.Fatalf("test %d: expected %v, got %v", i, tt.kind, tok.Kind)
		}

		if tok.Lexeme != tt.lexeme {
			t.Fatalf("test %d: expected %q, got %q", i, tt.lexeme, tok.Lexeme)
		}
	}
}

// func TestUnterminatedString(t *testing.T) {
// 	lexer := New(`"hello world`, "bad_string.flint")
// 	tok := lexer.Next()

// 	if tok.Kind != Illegal {
// 		t.Fatalf("expected Error token, got %v", tok.Kind)
// 	}
// }

// func TestInvalidByte(t *testing.T) {
// 	lexer := New(`'ab'`, "bad_byte.flint")
// 	tok := lexer.Next()

// 	if tok.Kind != Illegal {
// 		t.Fatalf("expected Error token, got %v", tok.Kind)
// 	}
// }

func TestEscapeSequences(t *testing.T) {
	input := `"\n\t\\\""`
	lexer := New(input, "escapes.flint")
	tok := lexer.Next()

	if tok.Kind != String {
		t.Fatalf("expected String, got %v", tok.Kind)
	}

	if tok.Lexeme != `"\n\t\\\""` {
		t.Fatalf("unexpected lexeme: %q", tok.Lexeme)
	}
}

func TestWhitespaceOnly(t *testing.T) {
	lexer := New("   \n\t   ", "space.flint")
	tok := lexer.Next()

	if tok.Kind != EndOfFile {
		t.Fatalf("expected EOF, got %v", tok.Kind)
	}
}

func TestEOFRepeatable(t *testing.T) {
	lexer := New("", "eof.flint")

	for i := 0; i < 5; i++ {
		tok := lexer.Next()
		if tok.Kind != EndOfFile {
			t.Fatalf("expected EOF on call %d, got %v", i, tok.Kind)
		}
	}
}
