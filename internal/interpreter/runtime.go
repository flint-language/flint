package interpreter

type ValueKind int

const (
	ValInt ValueKind = iota
	ValFloat
	ValUnsigned
	ValString
	ValBool
	ValNil
)

type Value struct {
	Kind     ValueKind
	Int      int64
	Float    float64
	Unsigned uint64
	String   string
	Bool     bool
}

func Int(v int64) Value {
	return Value{Kind: ValInt, Int: v}
}

func Float(v float64) Value {
	return Value{Kind: ValFloat, Float: v}
}

func Unsigned(v uint64) Value {
	return Value{Kind: ValUnsigned, Unsigned: v}
}

func String(v string) Value {
	return Value{Kind: ValString, String: v}
}

func Bool(v bool) Value {
	return Value{Kind: ValBool, Bool: v}
}

func Nil() Value {
	return Value{Kind: ValNil}
}
