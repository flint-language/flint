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
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyI8}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyI16}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyI32}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyI64}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyU8}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyU16}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyU32}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyU64}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Minus: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyI8}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyI16}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyI32}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyI64}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyU8}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyU16}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyU32}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyU64}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Star: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyI8}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyI16}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyI32}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyI64}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyU8}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyU16}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyU32}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyU64}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Slash: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyI8}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyI16}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyI32}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyI64}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyU8}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyU16}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyU32}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyU64}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.Percent: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyI8}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyI16}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyI32}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyI64}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyInt}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyU8}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyU16}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyU32}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyU64}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}},
	},
	lexer.PlusDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyF32}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyF64}},
	},
	lexer.MinusDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyF32}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyF64}},
	},
	lexer.StarDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyF32}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyF64}},
	},
	lexer.SlashDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyFloat}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyF32}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyF64}},
	},
	lexer.AmperAmper: {{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}}},
	lexer.VbarVbar:   {{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}}},
	lexer.Less: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyBool}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyBool}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyBool}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyBool}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyBool}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyBool}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyBool}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyBool}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyBool}},
	},
	lexer.LessDot: {
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyBool}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyBool}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyBool}},
	},
	lexer.EqualEqual: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyBool}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyBool}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyBool}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyBool}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyBool}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyBool}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyBool}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyBool}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyBool}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyBool}},
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyBool}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyBool}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyBool}},
		{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}},
		{Type{TKind: TyString}, Type{TKind: TyString}, Type{TKind: TyBool}},
		{Type{TKind: TyByte}, Type{TKind: TyByte}, Type{TKind: TyBool}},
	},
	lexer.NotEqual: {
		{Type{TKind: TyI8}, Type{TKind: TyI8}, Type{TKind: TyBool}},
		{Type{TKind: TyI16}, Type{TKind: TyI16}, Type{TKind: TyBool}},
		{Type{TKind: TyI32}, Type{TKind: TyI32}, Type{TKind: TyBool}},
		{Type{TKind: TyI64}, Type{TKind: TyI64}, Type{TKind: TyBool}},
		{Type{TKind: TyInt}, Type{TKind: TyInt}, Type{TKind: TyBool}},
		{Type{TKind: TyU8}, Type{TKind: TyU8}, Type{TKind: TyBool}},
		{Type{TKind: TyU16}, Type{TKind: TyU16}, Type{TKind: TyBool}},
		{Type{TKind: TyU32}, Type{TKind: TyU32}, Type{TKind: TyBool}},
		{Type{TKind: TyU64}, Type{TKind: TyU64}, Type{TKind: TyBool}},
		{Type{TKind: TyUnsigned}, Type{TKind: TyUnsigned}, Type{TKind: TyBool}},
		{Type{TKind: TyFloat}, Type{TKind: TyFloat}, Type{TKind: TyBool}},
		{Type{TKind: TyF32}, Type{TKind: TyF32}, Type{TKind: TyBool}},
		{Type{TKind: TyF64}, Type{TKind: TyF64}, Type{TKind: TyBool}},
		{Type{TKind: TyBool}, Type{TKind: TyBool}, Type{TKind: TyBool}},
		{Type{TKind: TyString}, Type{TKind: TyString}, Type{TKind: TyBool}},
		{Type{TKind: TyByte}, Type{TKind: TyByte}, Type{TKind: TyBool}},
	},
}

var unaryOps = map[lexer.TokenKind][]UnaryOpSig{
	lexer.Minus: {
		{Arg: Type{TKind: TyI8}, Out: Type{TKind: TyI8}},
		{Arg: Type{TKind: TyI16}, Out: Type{TKind: TyI16}},
		{Arg: Type{TKind: TyI32}, Out: Type{TKind: TyI32}},
		{Arg: Type{TKind: TyI64}, Out: Type{TKind: TyI64}},
		{Arg: Type{TKind: TyInt}, Out: Type{TKind: TyInt}},
		{Arg: Type{TKind: TyU8}, Out: Type{TKind: TyU8}},
		{Arg: Type{TKind: TyU16}, Out: Type{TKind: TyU16}},
		{Arg: Type{TKind: TyU32}, Out: Type{TKind: TyU32}},
		{Arg: Type{TKind: TyU64}, Out: Type{TKind: TyU64}},
		{Arg: Type{TKind: TyUnsigned}, Out: Type{TKind: TyUnsigned}},
	},
	lexer.MinusDot: {
		{Arg: Type{TKind: TyFloat}, Out: Type{TKind: TyFloat}},
		{Arg: Type{TKind: TyF32}, Out: Type{TKind: TyF32}},
		{Arg: Type{TKind: TyF64}, Out: Type{TKind: TyF64}},
	},
	lexer.Bang: {
		{Arg: Type{TKind: TyBool}, Out: Type{TKind: TyBool}},
	},
}
