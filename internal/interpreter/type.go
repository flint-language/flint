package interpreter

type TypeKind int

const (
	TyInt TypeKind = iota
	TyFloat
	TyUnsigned
	TyString
	TyBool
	TyUnit
)

type Type struct {
	Kind TypeKind
}

func (t Type) String() string {
	switch t.Kind {
	case TyInt:
		return "int"
	case TyFloat:
		return "float"
	case TyUnsigned:
		return "uint"
	case TyString:
		return "string"
	case TyBool:
		return "bool"
	case TyUnit:
		return "unit"
	default:
		return "<unknown>"
	}
}
