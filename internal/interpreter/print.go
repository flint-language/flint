package interpreter

import "fmt"

func printValue(v Value) string {
	switch v.Kind {
	case ValInt:
		return fmt.Sprintf("%d", v.Int)
	case ValFloat:
		return fmt.Sprintf("%g", v.Float)
	case ValUnsigned:
		return fmt.Sprintf("%d", v.Unsigned)
	case ValString:
		return fmt.Sprintf("%q", v.String)
	case ValBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case ValNil:
		return "()"
	default:
		return "<unknown>"
	}
}

func PrintReplResult(v Value) string {
	ty := TypeOf(v)
	return fmt.Sprintf("- : %s = %s", ty, printValue(v))
}
