package typechecker

import "flint/internal/lexer"

type BinOpSig struct {
	Left  Type
	Right Type
	Out   Type
}

type UnaryOpSig struct {
	Arg Type
	Out Type
}

var binOps = map[lexer.TokenKind][]BinOpSig{
	lexer.Plus: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Minus: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Star: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Slash: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Percent: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.PlusDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
	},
	lexer.MinusDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
	},
	lexer.StarDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
	},
	lexer.SlashDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
	},
	lexer.AmperAmper: {{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}}},
	lexer.VbarVbar:   {{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}}},
	lexer.Less: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyBool}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyBool}},
	},
	lexer.LessDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyBool}},
	},
	lexer.EqualEqual: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyBool}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyBool}},
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyBool}},
		{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}},
		{Type{TKind: TyString}, Type{TKind: TyString}, Type{TKind: TyBool}},
		{Type{TKind: TyByte}, Type{TKind: TyByte}, Type{TKind: TyBool}},
	},
	lexer.NotEqual: {
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyBool}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyBool}},
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyBool}},
		{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}},
		{Type{TKind: TyString}, Type{TKind: TyString}, Type{TKind: TyBool}},
		{Type{TKind: TyByte}, Type{TKind: TyByte}, Type{TKind: TyBool}},
	},
}

var unaryOps = map[lexer.TokenKind][]UnaryOpSig{
	lexer.Minus: {
		{Arg: Type{TKind: TyInt}, Out: Type{TKind: TyInt}},
		{Arg: Type{TKind: TyUnsigned}, Out: Type{TKind: TyUnsigned}},
	},
	lexer.MinusDot: {
		{Arg: Type{TKind: TyFloat}, Out: Type{TKind: TyFloat}},
	},
	lexer.Bang: {
		{Arg: Type{TKind: TyBool}, Out: Type{TKind: TyBool}},
	},
}
