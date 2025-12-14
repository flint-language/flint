package interpreter

func TypeOf(v Value) Type {
	switch v.Kind {
	case ValInt:
		return Type{Kind: TyInt}
	case ValFloat:
		return Type{Kind: TyFloat}
	case ValUnsigned:
		return Type{Kind: TyUnsigned}
	case ValString:
		return Type{Kind: TyString}
	case ValBool:
		return Type{Kind: TyBool}
	case ValNil:
		return Type{Kind: TyUnit}
	default:
		return Type{Kind: TyUnit}
	}
}
