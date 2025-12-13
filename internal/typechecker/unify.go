package typechecker

import "fmt"

type Subst map[int]*Type

func NewSubst() Subst {
	return Subst{}
}

func (s Subst) Follow(t *Type) *Type {
	if t == nil {
		return nil
	}
	for {
		if t.TKind == TyVar {
			if b, ok := s[t.VarID]; ok {
				t = b
				continue
			}
		}
		return t
	}
}

func (s Subst) Apply(t *Type) *Type {
	t = s.Follow(t)
	if t == nil {
		return nil
	}
	switch t.TKind {
	case TyVar:
		return t
	case TyInt, TyFloat, TyUnsigned, TyBool, TyString, TyByte, TyNil, TyError:
		return t
	case TyFunc:
		params := make([]*Type, len(t.Params))
		for i := range t.Params {
			params[i] = s.Apply(t.Params[i])
		}
		return &Type{
			TKind:  TyFunc,
			Params: params,
			Ret:    s.Apply(t.Ret),
		}
	case TyList:
		return &Type{
			TKind: TyList,
			Elem:  s.Apply(t.Elem),
		}
	case TyTuple:
		elems := make([]*Type, len(t.TElems))
		for i := range t.TElems {
			elems[i] = s.Apply(t.TElems[i])
		}
		return &Type{
			TKind:  TyTuple,
			TElems: elems,
		}
	default:
		return t
	}
}

func (s Subst) occurs(id int, t *Type) bool {
	t = s.Follow(t)
	if t == nil {
		return false
	}
	if t.TKind == TyVar {
		return t.VarID == id
	}
	switch t.TKind {
	case TyFunc:
		for _, p := range t.Params {
			if s.occurs(id, p) {
				return true
			}
		}
		return s.occurs(id, t.Ret)
	case TyList:
		return s.occurs(id, t.Elem)
	case TyTuple:
		for _, e := range t.TElems {
			if s.occurs(id, e) {
				return true
			}
		}
	}
	return false
}

func (s Subst) bindVar(v *Type, t *Type) error {
	v = s.Follow(v)
	t = s.Follow(t)
	if v == t {
		return nil
	}
	if v.TKind != TyVar {
		return fmt.Errorf("bindVar: left side is not a type variable: %v", v)
	}
	if s.occurs(v.VarID, t) {
		return fmt.Errorf("infinite type: %v occurs in %v", v, t)
	}
	if t.TKind == TyConcrete {
		v.Concrete = t.Concrete
		v.Family = t.Family
	}
	if t.TKind == TyVar && t.Concrete != CnUnknown {
		v.Concrete = t.Concrete
		v.Family = t.Family
	}
	if v.Family != FamilyUnknown && t.Family != FamilyUnknown && v.Family != t.Family {
		return fmt.Errorf("cannot unify type variables with different family constraints: %v vs %v", v.Family, t.Family)
	}

	s[v.VarID] = t
	return nil
}

func (s Subst) Unify(a, b *Type) error {
	a = s.Follow(a)
	b = s.Follow(b)
	if a == nil || b == nil {
		return fmt.Errorf("unify: nil type")
	}
	if a == b {
		return nil
	}
	if a.TKind == TyVar {
		return s.bindVar(a, b)
	}
	if b.TKind == TyVar {
		return s.bindVar(b, a)
	}
	if (a.TKind == TyInt || a.TKind == TyFloat || a.TKind == TyUnsigned) &&
		(b.TKind == TyInt || b.TKind == TyFloat || b.TKind == TyUnsigned) {
		if a.Concrete != CnUnknown && b.Concrete == CnUnknown {
			b.Concrete = a.Concrete
			b.Family = a.Family
		} else if b.Concrete != CnUnknown && a.Concrete == CnUnknown {
			a.Concrete = b.Concrete
			a.Family = b.Family
		}
		if a.Concrete != b.Concrete {
			return fmt.Errorf("cannot unify concrete types %s and %s", a.String(), b.String())
		}
		return nil
	}
	if a.TKind == TyFunc && b.TKind == TyFunc {
		if len(a.Params) != len(b.Params) {
			return fmt.Errorf("function arity mismatch: %d vs %d", len(a.Params), len(b.Params))
		}
		for i := 0; i < len(a.Params); i++ {
			if err := s.Unify(a.Params[i], b.Params[i]); err != nil {
				return err
			}
		}
		return s.Unify(a.Ret, b.Ret)
	}
	if a.TKind == TyList && b.TKind == TyList {
		return s.Unify(a.Elem, b.Elem)
	}
	if a.TKind == TyTuple && b.TKind == TyTuple {
		if len(a.TElems) != len(b.TElems) {
			return fmt.Errorf("tuple length mismatch: %d vs %d", len(a.TElems), len(b.TElems))
		}
		for i := range a.TElems {
			if err := s.Unify(a.TElems[i], b.TElems[i]); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("cannot unify types %s and %s", a.String(), b.String())
}

func (t *Type) Fresh() *Type {
	if t == nil {
		return &Type{TKind: TyError}
	}
	switch t.TKind {
	case TyVar, TyInt, TyFloat, TyUnsigned, TyBool, TyByte, TyString:
		return &Type{
			TKind:    t.TKind,
			VarID:    t.VarID,
			Family:   t.Family,
			Concrete: t.Concrete,
			Elem:     t.Elem,
			TElems:   t.TElems,
			Params:   t.Params,
			Ret:      t.Ret,
		}
	case TyFunc:
		newParams := make([]*Type, len(t.Params))
		for i, p := range t.Params {
			newParams[i] = p.Fresh()
		}
		return &Type{
			TKind:  TyFunc,
			Params: newParams,
			Ret:    t.Ret.Fresh(),
		}
	case TyList:
		var elem *Type
		if t.Elem != nil {
			elem = t.Elem.Fresh()
		}
		return &Type{TKind: TyList, Elem: elem}
	case TyTuple:
		newElems := make([]*Type, len(t.TElems))
		for i, e := range t.TElems {
			newElems[i] = e.Fresh()
		}
		return &Type{TKind: TyTuple, TElems: newElems}
	case TyRange:
		var elem *Type
		if t.Elem != nil {
			elem = t.Elem.Fresh()
		}
		return &Type{TKind: TyRange, Elem: elem}
	default:
		return &Type{TKind: TyError}
	}
}
