package runtime

func (v Value) AsInt() int64 {
	if v.Kind != ValInt {
		TypeError("expected int")
	}
	return v.Int
}

func (v Value) AsFloat() float64 {
	if v.Kind != ValFloat {
		TypeError("expected float")
	}
	return v.Float
}

func (v Value) AsBool() bool {
	if v.Kind != ValBool {
		TypeError("expected bool")
	}
	return v.Bool
}
