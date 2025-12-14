package vm

import (
	"flint/internal/runtime"
)


type Stack struct {
	data []runtime.Value
}

func NewStack() Stack {
	return Stack{data: make([]runtime.Value, 0, 256)}
}

func (s *Stack) Push(v runtime.Value) {
	s.data = append(s.data, v)
}

func (s *Stack) Pop() runtime.Value {
	n := len(s.data)
	if n == 0 {
		runtime.RuntimeError("stack underflow")
	}
	v := s.data[n-1]
	s.data = s.data[:n-1]
	return v
}

func (s *Stack) Peek() runtime.Value {
	return s.data[len(s.data)-1]
}
