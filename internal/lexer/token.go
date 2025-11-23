package lexer

// TokenKind is an enum-like type describing the category of a token.
// Using a custom type instead of strings avoids mistakes and improves performance.
type TokenKind int

// Token represents a lexical token: kind, literal text, and source position.
type Token struct {
	Kind   TokenKind // token category
	Lexeme string    // exact source text
	Line   int       // 1-based line number
	Column int       // 1-based column number
}

// Token kinds for the Flint language.
// These are grouped for clarity and to ease future extension.
const (
	// Special
	Illegal TokenKind = iota
	Comment

	// Identifiers and literals
	Identifier
	Int
	Float
	String
	Byte
	Bool
	Tuple
	List
	Nil

	// Groupings
	LeftParen    // (
	RightParen   // )
	LeftBrace    // {
	RightBrace   // }
	LeftBracket  // [
	RightBracket // ]

	// Int Operations
	Plus
	Minus
	Star
	Slash
	Less
	Greater
	LessEqual
	GreaterEqual
	Percent

	// Float Operations
	PlusDot         // '+.'
	MinusDot        // '-.'
	StarDot         // '*.'
	SlashDot        // '/.'
	LessDot         // '<.'
	GreaterDot      // '>.'
	LessEqualDot    // '<=.'
	GreaterEqualDot // '>=.'

	// String Operation
	LtGt // '<>'

	// Other Punctuation
	Colon
	Comma
	Bang // '!'
	Equal
	EqualEqual // '=='
	NotEqual   // '!='
	Vbar       // '|'
	VbarVbar   // '||'
	AmperAmper // '&&'
	Pipe       // '|>'
	Dot        // '.'
	RArrow     // '->'
	DotDot     // '..'
	At         // '@'
	Underscore // '_'
	EndOfFile

	// Keywords
	KwAs
	KwAssert
	KwBool
	KwByte
	KwElse
	KwFloat
	KwFn
	KwFor
	KwIf
	KwIn
	KwInt
	KwList
	KwMatch
	KwMut
	KwNil
	KwPanic
	KwPub
	KwString
	KwType
	KwUse
	KwVal
	KwWhere
)

var KeywordMap = map[string]TokenKind{
	"as":     KwAs,
	"assert": KwAssert,
	"Bool":   KwBool,
	"Byte":   KwByte,
	"else":   KwElse,
	"False":  Bool,
	"Float":  KwFloat,
	"fn":     KwFn,
	"for":    KwFor,
	"if":     KwIf,
	"in":     KwIn,
	"Int":    KwInt,
	"List":   KwList,
	"match":  KwMatch,
	"mut":    KwMut,
	"Nil":    KwNil,
	"panic":  KwPanic,
	"pub":    KwPub,
	"String": KwString,
	"True":   Bool,
	"type":   KwType,
	"use":    KwUse,
	"val":    KwVal,
	"where":  KwWhere,
}

func LookupIdentifier(name string) TokenKind {
	if k, ok := KeywordMap[name]; ok {
		return k
	}
	return Identifier
}

var precedence = map[TokenKind]int{
	VbarVbar:   1,
	Pipe:       1,
	Colon:      1,
	AmperAmper: 2,
	EqualEqual: 3,
	NotEqual:   3,

	Less:            4,
	LessEqual:       4,
	LessDot:         4,
	LessEqualDot:    4,
	Greater:         4,
	GreaterEqual:    4,
	GreaterDot:      4,
	GreaterEqualDot: 4,

	Plus:     5,
	PlusDot:  5,
	Minus:    5,
	MinusDot: 5,
	LtGt:     5,

	Star:     6,
	StarDot:  6,
	Slash:    6,
	SlashDot: 6,
	Percent:  6,
	DotDot:   6,
}

func (k TokenKind) Precedence() int {
	if p, ok := precedence[k]; ok {
		return p
	}
	return 0
}
