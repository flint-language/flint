package lexer

// TokenKind is an enum-like type describing the category of a token.
// Using a custom type instead of strings avoids mistakes and improves performance.
type TokenKind int

// Token represents a lexical token: kind, literal text, and source position.
type Token struct {
	Kind   TokenKind
	Lexeme string
	Line   int
	Column int

	File   string
	Source []rune
}

const (
	Illegal TokenKind = iota
	Comment

	Identifier
	Int
	Float
	String
	Byte
	Bool
	Tuple
	List
	Nil

	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	Plus
	Minus
	Star
	Slash
	Less
	Greater
	LessEqual
	GreaterEqual
	Percent

	PlusDot
	MinusDot
	StarDot
	SlashDot
	LessDot
	GreaterDot
	LessEqualDot
	GreaterEqualDot

	LtGt

	Colon
	Comma
	Bang
	Equal
	EqualEqual
	NotEqual
	Vbar
	VbarVbar
	AmperAmper
	Pipe
	Dot
	RArrow
	DotDot
	At
	Underscore
	EndOfFile

	KwAs
	KwBool
	KwByte
	KwElse
	KwFloat
	KwFn
	KwIf
	KwIn
	KwInt
	KwList
	KwMatch
	KwMut
	KwNil
	KwPub
	KwString
	KwThen
	KwType
	KwUse
	KwVal
)

var KeywordMap = map[string]TokenKind{
	"as":     KwAs,
	"Bool":   KwBool,
	"Byte":   KwByte,
	"else":   KwElse,
	"False":  Bool,
	"Float":  KwFloat,
	"fn":     KwFn,
	"if":     KwIf,
	"Int":    KwInt,
	"List":   KwList,
	"match":  KwMatch,
	"mut":    KwMut,
	"Nil":    KwNil,
	"pub":    KwPub,
	"String": KwString,
	"then":   KwThen,
	"True":   Bool,
	"type":   KwType,
	"use":    KwUse,
	"val":    KwVal,
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
