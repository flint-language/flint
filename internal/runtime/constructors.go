package runtime

func Int(v int64) Value {
	return Value{Kind: ValInt, Int: v}
}

func Float(v float64) Value {
	return Value{Kind: ValFloat, Float: v}
}

func Bool(v bool) Value {
	return Value{Kind: ValBool, Bool: v}
}

func String(v string) Value {
	return Value{Kind: ValString, Str: v}
}

func Stringify(v Value) string { return v.Str }

func (a Value) Equals(b Value) bool {
	if a.Kind != b.Kind {
		return false
	}

	switch a.Kind {
	case ValInt:
		return a.Int == b.Int
	case ValFloat:
		return a.Float == b.Float
	case ValBool:
		return a.Bool == b.Bool
	case ValString:
		return a.Str == b.Str
	default:
		TypeError("equality not supported")
	}
	return false
}

func (v Value) IsTruthy() bool {
	switch v.Kind {
	case ValBool:
		return v.Bool
	case ValInt:
		return v.Int != 0
	case ValFloat:
		return v.Float != 0.0
	case ValString:
		return v.Str != ""
	default:
		return false
	}
}
