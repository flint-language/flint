package typechecker

import "flint/internal/parser"

var PlatformIntBits int = 0

func (tc *TypeChecker) resolveType(t parser.Expr) *Type {
	switch typ := t.(type) {
	case *parser.TypeExpr:
		switch typ.Name {
		case "Int":
			return NewTypeVar(FamilyInt)
		case "Float":
			return NewTypeVar(FamilyFloat)
		case "Unsigned":
			return NewTypeVar(FamilyUnsigned)
		case "I8":
			return &Type{TKind: TyConcrete, Concrete: CI8, Family: FamilyInt}
		case "I16":
			return &Type{TKind: TyConcrete, Concrete: CI16, Family: FamilyInt}
		case "I32":
			return &Type{TKind: TyConcrete, Concrete: CI32, Family: FamilyInt}
		case "I64":
			return &Type{TKind: TyConcrete, Concrete: CI64, Family: FamilyInt}
		case "U8":
			return &Type{TKind: TyConcrete, Concrete: CU8, Family: FamilyUnsigned}
		case "U16":
			return &Type{TKind: TyConcrete, Concrete: CU16, Family: FamilyUnsigned}
		case "U32":
			return &Type{TKind: TyConcrete, Concrete: CU32, Family: FamilyUnsigned}
		case "U64":
			return &Type{TKind: TyConcrete, Concrete: CU64, Family: FamilyUnsigned}
		case "F32":
			return &Type{TKind: TyConcrete, Concrete: CF32, Family: FamilyFloat}
		case "F64":
			return &Type{TKind: TyConcrete, Concrete: CF64, Family: FamilyFloat}
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

func isLiteral(t *Type) bool {
	switch t.TKind {
	case TyVar:
		return t.Family == FamilyInt || t.Family == FamilyFloat || t.Family == FamilyUnsigned
	case TyInt, TyUnsigned, TyFloat:
		return true
	default:
		return false
	}
}

func sameFamily(a, b *Type) bool {
	fa := a.Family
	fb := b.Family
	if a.TKind == TyInt || a.TKind == TyUnsigned || a.TKind == TyFloat {
		fa = a.Family
	}
	if b.TKind == TyInt || b.TKind == TyUnsigned || b.TKind == TyFloat {
		fb = b.Family
	}
	return fa == fb
}
