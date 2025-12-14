package runtime

func RuntimeError(msg string) {
	panic("runtime error: " + msg)
}

func RuntimeErrorNoMsg() {
	panic("runtime error")
}

func TypeError(msg string) {
	panic("type error: " + msg)
}

func TypeErrorNoMsg() {
	panic("type error")
}

func MathError(msg string) {
	panic("math error: " + msg)
}

func MathErrorNoMsg() {
	panic("math error")
}
