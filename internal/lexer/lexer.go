package lexer

import (
	"fmt"
	"strconv"
	"unicode"
)

type Lexer struct {
	source       []rune
	position     int
	lineNumber   int
	columnNumber int
	fileName     string
}

func Tokenize(source, filename string) ([]Token, error) {
	lx := New(source, filename)
	out := []Token{}
	for {
		tok := lx.Next()
		out = append(out, tok)
		if tok.Kind == EndOfFile {
			break
		}
		if tok.Kind == Illegal {
			lx.error("Illegal character")
		}
	}
	return out, nil
}

func New(source, filename string) *Lexer {
	r := []rune(source)
	return &Lexer{source: r, position: 0, lineNumber: 1, columnNumber: 1, fileName: filename}
}

func (l *Lexer) Next() Token {
	l.consumeWhitespace()
	startlineNumber, startcolumnNumber := l.lineNumber, l.columnNumber
	ch := l.peekRuneAt(0)
	if ch == 0 {
		return l.makeToken(EndOfFile, "", startlineNumber, startcolumnNumber)
	}
	if isIdentifierStart(ch) {
		lex := l.scanIdentifier()
		kind := LookupIdentifier(lex)
		return l.makeToken(kind, lex, startlineNumber, startcolumnNumber)
	}
	if unicode.IsDigit(ch) {
		lex, kind := l.scanNumberLiteral()
		return l.makeToken(kind, lex, startlineNumber, startcolumnNumber)
	}
	if ch == '"' {
		return l.makeToken(String, l.scanStringLiteral(), startlineNumber, startcolumnNumber)
	}
	if ch == '\'' {
		return l.makeToken(Byte, l.scanByteLiteral(), startlineNumber, startcolumnNumber)
	}
	if ch == '/' && l.peekRuneAt(1) == '*' {
		return l.scanBlockComment()
	}
	if ch == '/' && l.peekRuneAt(1) == '/' {
		return l.makeToken(Comment, l.scanLineComment(), startlineNumber, startcolumnNumber)
	}
	switch ch {
	case '=':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			return l.makeToken(EqualEqual, "==", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Equal, "=", startlineNumber, startcolumnNumber)

	case '!':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			return l.makeToken(NotEqual, "!=", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Bang, "!", startlineNumber, startcolumnNumber)
	case '<':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			if l.peekRuneAt(0) == '.' {
				l.advanceRune()
				return l.makeToken(LessEqualDot, "<=.", startlineNumber, startcolumnNumber)
			}
			return l.makeToken(LessEqual, "<=", startlineNumber, startcolumnNumber)
		}
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(LessDot, "<.", startlineNumber, startcolumnNumber)
		}
		if l.peekRuneAt(0) == '>' {
			l.advanceRune()
			return l.makeToken(LtGt, "<>", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Less, "<", startlineNumber, startcolumnNumber)
	case '>':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			if l.peekRuneAt(0) == '.' {
				l.advanceRune()
				return l.makeToken(GreaterEqualDot, ">=.", startlineNumber, startcolumnNumber)
			}
			return l.makeToken(GreaterEqual, ">=", startlineNumber, startcolumnNumber)
		}
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(GreaterDot, ">.", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Greater, ">", startlineNumber, startcolumnNumber)
	case '+':
		l.advanceRune()
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(PlusDot, "+.", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Plus, "+", startlineNumber, startcolumnNumber)
	case '-':
		l.advanceRune()
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(MinusDot, "-.", startlineNumber, startcolumnNumber)
		}
		if l.peekRuneAt(0) == '>' {
			l.advanceRune()
			return l.makeToken(RArrow, "->", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Minus, "-", startlineNumber, startcolumnNumber)
	case '*':
		l.advanceRune()
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(StarDot, "*.", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Star, "*", startlineNumber, startcolumnNumber)
	case '/':
		l.advanceRune()
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(SlashDot, "/.", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Slash, "/", startlineNumber, startcolumnNumber)
	case '%':
		l.advanceRune()
		return l.makeToken(Percent, "%", startlineNumber, startcolumnNumber)
	case ':':
		l.advanceRune()
		return l.makeToken(Colon, ":", startlineNumber, startcolumnNumber)
	case ',':
		l.advanceRune()
		return l.makeToken(Comma, ",", startlineNumber, startcolumnNumber)
	case '{':
		l.advanceRune()
		return l.makeToken(LeftBrace, "{", startlineNumber, startcolumnNumber)
	case '}':
		l.advanceRune()
		return l.makeToken(RightBrace, "}", startlineNumber, startcolumnNumber)
	case '(':
		l.advanceRune()
		return l.makeToken(LeftParen, "(", startlineNumber, startcolumnNumber)
	case ')':
		l.advanceRune()
		return l.makeToken(RightParen, ")", startlineNumber, startcolumnNumber)
	case '[':
		l.advanceRune()
		return l.makeToken(LeftBracket, "[", startlineNumber, startcolumnNumber)
	case ']':
		l.advanceRune()
		return l.makeToken(RightBracket, "]", startlineNumber, startcolumnNumber)
	case '|':
		l.advanceRune()
		if l.peekRuneAt(0) == '|' {
			l.advanceRune()
			return l.makeToken(VbarVbar, "||", startlineNumber, startcolumnNumber)
		}
		if l.peekRuneAt(0) == '>' {
			l.advanceRune()
			return l.makeToken(Pipe, "|>", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Vbar, "|", startlineNumber, startcolumnNumber)
	case '&':
		l.advanceRune()
		if l.peekRuneAt(0) == '&' {
			l.advanceRune()
			return l.makeToken(AmperAmper, "&&", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Illegal, "&", startlineNumber, startcolumnNumber)
	case '.':
		l.advanceRune()
		if l.peekRuneAt(0) == '.' {
			l.advanceRune()
			return l.makeToken(DotDot, "..", startlineNumber, startcolumnNumber)
		}
		return l.makeToken(Dot, ".", startlineNumber, startcolumnNumber)
	case '@':
		l.advanceRune()
		return l.makeToken(At, "@", startlineNumber, startcolumnNumber)
	default:
		return l.makeToken(Illegal, string(l.advanceRune()), startlineNumber, startcolumnNumber)
	}
}

func (l *Lexer) advanceRune() rune {
	if l.position >= len(l.source) {
		return 0
	}
	r := l.source[l.position]
	l.position++
	if r == '\n' {
		l.lineNumber++
		l.columnNumber = 1
	} else {
		l.columnNumber++
	}
	return r
}

func (l *Lexer) peekRuneAt(offset int) rune {
	idx := l.position + offset
	if idx >= len(l.source) {
		return 0
	}
	return l.source[idx]
}

func (l *Lexer) makeToken(kind TokenKind, lexeme string, lineNumber, columnNumber int) Token {
	return Token{
		Kind:   kind,
		Lexeme: lexeme,
		Line:   lineNumber,
		Column: columnNumber,
		File:   l.fileName,
		Source: l.source,
	}
}

func (l *Lexer) scanIdentifier() string {
	start := l.position
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 || !(isIdentifierPart(ch)) {
			break
		}
		l.advanceRune()
	}
	return string(l.source[start:l.position])
}

func (l *Lexer) scanNumberLiteral() (string, TokenKind) {
	start := l.position
	isFloat := false
	for {
		ch := l.peekRuneAt(0)
		if ch == '.' {
			if l.peekRuneAt(1) == '.' {
				break
			}
			isFloat = true
			l.advanceRune()
			continue
		}
		if !unicode.IsDigit(ch) && ch != '_' {
			break
		}
		l.advanceRune()
	}
	if l.peekRuneAt(0) == 'u' {
		l.advanceRune()
		lex := string(l.source[start:l.position])
		clean := StripNumericSeparators(lex[:len(lex)-1])
		if _, err := strconv.ParseUint(clean, 10, 64); err == nil {
			return lex, Unsigned
		}
		return lex, Illegal
	}
	lex := string(l.source[start:l.position])
	clean := StripNumericSeparators(lex)
	if isFloat {
		if _, err := strconv.ParseFloat(clean, 64); err == nil {
			return lex, Float
		}
	}
	if _, err := strconv.ParseInt(clean, 10, 64); err == nil {
		return lex, Int
	}
	if _, err := strconv.ParseFloat(clean, 64); err == nil {
		return lex, Float
	}
	return lex, Illegal
}

func (l *Lexer) scanStringLiteral() string {
	quote := l.advanceRune()
	start := l.position - 1
	if quote != '"' {
		return string(l.source[start:l.position])
	}
	ch := l.advanceRune()
	if ch == 0 {
		l.error("unterminated string literal")
		return string(l.source[start:l.position])
	}
	if ch == quote {
		l.error("empty string literal")
		return string(l.source[start:l.position])
	}
	runeCount := 0
	for {
		if ch == '\\' {
			runeCount++
			esc := l.advanceRune()
			if esc == 0 {
				l.error("unterminated escape sequence in string literal")
				return string(l.source[start:l.position])
			}
			switch esc {
			case 'n', 't', 'r', '\\', '\'', '"', '0':
			default:
				l.error(fmt.Sprintf("invalid escape character: \\%c", esc))
				return string(l.source[start:l.position])
			}
		} else {
			runeCount++
		}
		ch = l.advanceRune()
		if ch == 0 {
			l.error("unterminated string literal")
			return string(l.source[start:l.position])
		}
		if ch == quote {
			break
		}
	}
	return string(l.source[start:l.position])
}

func (l *Lexer) scanByteLiteral() string {
	quote := l.advanceRune()
	start := l.position - 1
	if quote != '\'' {
		return string(l.source[start:l.position])
	}
	ch := l.advanceRune()
	if ch == 0 {
		l.error("unterminated character literal")
		return string(l.source[start:l.position])
	}
	if ch == quote {
		l.error("empty character literal")
		return string(l.source[start:l.position])
	}
	var runeCount int = 0
	if ch == '\\' {
		runeCount++
		esc := l.advanceRune()
		if esc == 0 {
			l.error("unterminated escape sequence in character literal")
			return string(l.source[start:l.position])
		}
		switch esc {
		case 'n', 't', 'r', '\\', '\'', '"', '0':
		default:
			l.error(fmt.Sprintf("invalid escape character: \\%c", esc))
			return string(l.source[start:l.position])
		}
	} else {
		runeCount++
	}
	end := l.advanceRune()
	if end == 0 {
		l.error("unterminated character literal")
		return string(l.source[start:l.position])
	}
	if end != '\'' {
		l.error("extra characters in character literal (expected closing ')")
		return string(l.source[start:l.position])
	}
	if runeCount != 1 {
		l.error("character literal must contain exactly 1 character or valid escape")
	}
	return string(l.source[start:l.position])
}

func (l *Lexer) scanLineComment() string {
	start := l.position
	l.advanceRune()
	l.advanceRune()
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 || ch == '\n' {
			break
		}
		l.advanceRune()
	}
	return string(l.source[start:l.position])
}

func (l *Lexer) scanBlockComment() Token {
	startLine, startCol := l.lineNumber, l.columnNumber
	start := l.position
	l.advanceRune()
	l.advanceRune()
	isDoc := false
	if l.peekRuneAt(0) == '*' {
		isDoc = true
		l.advanceRune()
	}
	for {
		ch := l.advanceRune()
		if ch == 0 {
			l.error("unterminated block comment")
			break
		}
		if ch == '*' && l.peekRuneAt(0) == '/' {
			l.advanceRune()
			break
		}
	}
	lexeme := string(l.source[start:l.position])
	kind := Comment
	if isDoc {
		kind = DocComment
	}
	return l.makeToken(kind, lexeme, startLine, startCol)
}

func (l *Lexer) consumeWhitespace() {
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 {
			return
		}
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advanceRune()
			continue
		}
		break
	}
}
