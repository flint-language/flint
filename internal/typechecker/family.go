package typechecker

import "fmt"

type FamilyKind int

const (
	FamilyUnknown FamilyKind = iota
	FamilyInt
	FamilyFloat
	FamilyUnsigned
)

func (f FamilyKind) String() string {
	switch f {
	case FamilyInt:
		return "Int"
	case FamilyFloat:
		return "Float"
	case FamilyUnsigned:
		return "Unsigned"
	default:
		return "Unknown"
	}
}

type ConcreteKind int

const (
	CnUnknown ConcreteKind = iota
	CI8
	CI16
	CI32
	CI64
	CF32
	CF64
	CU8
	CU16
	CU32
	CU64
)

func (c ConcreteKind) String() string {
	switch c {
	case CI8:
		return "I8"
	case CI16:
		return "I16"
	case CI32:
		return "I32"
	case CI64:
		return "I64"
	case CF32:
		return "F32"
	case CF64:
		return "F64"
	case CU8:
		return "U8"
	case CU16:
		return "U16"
	case CU32:
		return "U32"
	case CU64:
		return "U64"
	default:
		return "UnknownConcrete"
	}
}

type IType struct {
	Kind     TypeKind
	ID       int
	Family   FamilyKind
	Concrete ConcreteKind
	Param    *Type
	Result   *Type
}

func FamilyOfConcrete(c ConcreteKind) FamilyKind {
	switch c {
	case CI8, CI16, CI32, CI64:
		return FamilyInt
	case CF32, CF64:
		return FamilyFloat
	case CU8, CU16, CU32, CU64:
		return FamilyUnsigned
	default:
		return FamilyUnknown
	}
}

func ConcreteMembers(f FamilyKind) []ConcreteKind {
	switch f {
	case FamilyInt:
		return []ConcreteKind{CI8, CI16, CI32, CI64}
	case FamilyFloat:
		return []ConcreteKind{CF32, CF64}
	case FamilyUnsigned:
		return []ConcreteKind{CU8, CU16, CU32, CU64}
	default:
		return nil
	}
}

func DefaultConcreteForFamily(f FamilyKind) ConcreteKind {
	switch f {
	case FamilyInt:
		if PlatformIntBits == 32 {
			return CI32
		}
		return CI64
	case FamilyFloat:
		return CF64
	case FamilyUnsigned:
		if PlatformIntBits == 32 {
			return CU32
		}
		return CU64
	default:
		return CnUnknown
	}
}

func ConcreteToType(c ConcreteKind) *Type {
	switch c {
	case CI8, CI16, CI32, CI64:
		return &Type{TKind: TyInt, Concrete: c}
	case CF32, CF64:
		return &Type{TKind: TyFloat, Concrete: c}
	case CU8, CU16, CU32, CU64:
		return &Type{TKind: TyUnsigned, Concrete: c}
	default:
		return &Type{TKind: TyError}
	}
}

func ConcreteFromString(s string) (ConcreteKind, error) {
	switch s {
	case "I8":
		return CI8, nil
	case "I16":
		return CI16, nil
	case "I32":
		return CI32, nil
	case "I64":
		return CI64, nil
	case "F32":
		return CF32, nil
	case "F64":
		return CF64, nil
	case "U8":
		return CU8, nil
	case "U16":
		return CU16, nil
	case "U32":
		return CU32, nil
	case "U64":
		return CU64, nil
	default:
		return CnUnknown, fmt.Errorf("unknown concrete kind %q", s)
	}
}
