package typechecker

import (
	"fmt"
	"strconv"
	"strings"
)

type TypeKind int

type Type struct {
	TKind  TypeKind
	Params []*Type
	Ret    *Type
	Elem   *Type
	TElems []*Type

	VarID    int
	Family   FamilyKind
	Concrete ConcreteKind
}

const (
	TyError TypeKind = iota
	TyConcrete
	TyVar
	TyFamily
	TyInt
	TyFloat
	TyUnsigned
	TyBool
	TyByte
	TyString
	TyNil
	TyFunc
	TyList
	TyTuple
	TyRange
)

func (t Type) String() string {
	switch t.TKind {
	case TyVar:
		if t.Concrete != CnUnknown {
			return t.Concrete.String()
		}
		if t.Family != FamilyUnknown {
			return fmt.Sprintf("α%d:%s", t.VarID, t.Family.String())
		}
		return fmt.Sprintf("α%d", t.VarID)
	case TyInt:
		return "Int"
	case TyFloat:
		return "Float"
	case TyUnsigned:
		return "Unsigned"
	case TyBool:
		return "Bool"
	case TyByte:
		return "Byte"
	case TyString:
		return "String"
	case TyNil:
		return "Nil"
	case TyList:
		if t.Elem != nil {
			return fmt.Sprintf("List(%s)", t.Elem.String())
		}
		return "List(<unknown>)"
	case TyTuple:
		parts := []string{}
		for _, e := range t.TElems {
			if e == nil {
				parts = append(parts, "<unknown>")
			} else {
				parts = append(parts, e.String())
			}
		}
		return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
	case TyRange:
		if t.Elem != nil {
			return fmt.Sprintf("Range(%s)", t.Elem.String())
		}
		return "Range(Int)"
	case TyFunc:
		params := []string{}
		for _, p := range t.Params {
			params = append(params, p.String())
		}
		ret := "<unknown>"
		if t.Ret != nil {
			ret = t.Ret.String()
		}
		return fmt.Sprintf("(%s) -> %s", strings.Join(params, ", "), ret)
	case TyError:
		return "<error>"
	default:
		return fmt.Sprintf("<unknown:%d>", t.TKind)
	}
}

func (t Type) Kind() TypeKind { return t.TKind }

func (t *Type) Equal(u *Type) bool {
	if t == nil || u == nil {
		return t == u
	}
	if t.TKind != u.TKind {
		return false
	}
	switch t.TKind {
	case TyFunc:
		if len(t.Params) != len(u.Params) {
			return false
		}
		for i := range t.Params {
			if !t.Params[i].Equal(u.Params[i]) {
				return false
			}
		}
		return t.Ret.Equal(u.Ret)
	case TyList:
		if t.Elem == nil || u.Elem == nil {
			return t.Elem == u.Elem
		}
		return t.Elem.Equal(u.Elem)
	case TyTuple:
		if len(t.TElems) != len(u.TElems) {
			return false
		}
		for i := range t.TElems {
			if !t.TElems[i].Equal(u.TElems[i]) {
				return false
			}
		}
		return true
	case TyRange:
		if t.Elem == nil || u.Elem == nil {
			return t.Elem == u.Elem
		}
		return t.Elem.Equal(u.Elem)
	default:
		return true
	}
}

func init() {
	if strconv.IntSize == 32 {
		PlatformIntBits = 32
	} else {
		PlatformIntBits = 64
	}
}
