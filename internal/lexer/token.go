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
	DocComment

	Identifier
	Int
	Float
	Unsigned
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
	KwF32
	KwF64
	KwFn
	KwFun
	KwIf
	KwIn
	KwInt
	KwI8
	KwI16
	KwI32
	KwI64
	KwList
	KwMatch
	KwMut
	KwNil
	KwPub
	KwString
	KwThen
	KwType
	KwUse
	KwUnsigned
	KwU8
	KwU16
	KwU32
	KwU64
	KwVal
)

var KeywordMap = map[string]TokenKind{
	"as":       KwAs,
	"Bool":     KwBool,
	"Byte":     KwByte,
	"else":     KwElse,
	"False":    Bool,
	"Float":    KwFloat,
	"F32":      KwF32,
	"F64":      KwF64,
	"fn":       KwFn,
	"fun":      KwFun,
	"if":       KwIf,
	"Int":      KwInt,
	"I8":       KwI8,
	"I16":      KwI16,
	"I32":      KwI32,
	"I64":      KwI64,
	"List":     KwList,
	"match":    KwMatch,
	"mut":      KwMut,
	"Nil":      KwNil,
	"pub":      KwPub,
	"String":   KwString,
	"then":     KwThen,
	"True":     Bool,
	"type":     KwType,
	"use":      KwUse,
	"Unsigned": KwUnsigned,
	"U8":       KwU8,
	"U16":      KwU16,
	"U32":      KwU32,
	"U64":      KwU64,
	"val":      KwVal,
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
