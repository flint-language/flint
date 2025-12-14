package typechecker

import "flint/internal/parser"

var PlatformIntBits int = 0

func (tc *TypeChecker) resolveType(t parser.Expr) *Type {
	switch typ := t.(type) {
	case *parser.TypeExpr:
		switch typ.Name {
		case "Int":
			return &Type{TKind: TyInt}
		case "Float":
			return &Type{TKind: TyFloat}
		case "Unsigned":
			return &Type{TKind: TyUnsigned}
		case "I8":
			return &Type{TKind: TyI8}
		case "I16":
			return &Type{TKind: TyI16}
		case "I32":
			return &Type{TKind: TyI32}
		case "I64":
			return &Type{TKind: TyI64}
		case "U8":
			return &Type{TKind: TyU8}
		case "U16":
			return &Type{TKind: TyU16}
		case "U32":
			return &Type{TKind: TyU32}
		case "U64":
			return &Type{TKind: TyU64}
		case "F32":
			return &Type{TKind: TyF32}
		case "F64":
			return &Type{TKind: TyF64}
		case "Bool":
			return &Type{TKind: TyBool}
		case "String":
			return &Type{TKind: TyString}
		case "Byte":
			return &Type{TKind: TyByte}
		case "Nil":
			return &Type{TKind: TyNil}
		case "List":
			elemTy := &Type{TKind: TyNil}
			if typ.Generic != nil {
				elemTy = tc.resolveType(typ.Generic)
			}
			return &Type{TKind: TyList, Elem: elemTy}
		}
	case *parser.TupleTypeExpr:
		elems := []*Type{}
		for _, te := range typ.Types {
			elems = append(elems, tc.resolveType(te))
		}
		return &Type{TKind: TyTuple, TElems: elems}
	}

	return &Type{TKind: TyError}
}

func (e *Env) currentScopeGet(name string) (*VarInfo, bool) {
	ty, ok := e.vars[name]
	if !ok {
		return nil, false
	}
	return &ty, true
}

func (tc *TypeChecker) coerceLiteralTo(expected, actual *Type) *Type {
	if actual == nil || expected == nil {
		return actual
	}
	if actual.TKind == TyInt {
		switch expected.TKind {
		case TyI8, TyI16, TyI32, TyI64, TyInt:
			return expected
		}
	}
	if actual.TKind == TyUnsigned {
		switch expected.TKind {
		case TyU8, TyU16, TyU32, TyU64, TyUnsigned:
			return expected
		}
	}
	if actual.TKind == TyFloat {
		switch expected.TKind {
		case TyF32, TyF64, TyFloat:
			return expected
		}
	}
	return actual
}
