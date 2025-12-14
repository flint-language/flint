package runtime

import "fmt"

type ValueKind int

const (
	ValInt ValueKind = iota
	ValFloat
	ValBool
	ValString
)

type Value struct {
	Kind  ValueKind
	Int   int64
	Float float64
	Bool  bool
	Str   string
}

func (v Value) String() string {
	switch v.Kind {
	case ValInt:
		return fmt.Sprintf("%d", v.Int)

	case ValFloat:
		return fmt.Sprintf("%g", v.Float)

	case ValBool:
		return fmt.Sprintf("%t", v.Bool)

	case ValString:
		return fmt.Sprintf("%q", v.Str)

	// case Void:
	// 	return "nil"

	default:
		return v.Str
	}
}