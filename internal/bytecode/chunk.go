package bytecode

import (
	"flint/internal/runtime"
)

type Function struct {
	Name   string
	Chunk  *Chunk
	Params int
}

type Chunk struct {
	Code   []OpCode
	Consts []runtime.Value
	Funcs  []*Function
}

func (f *Function) String() string {
	return f.Name
}