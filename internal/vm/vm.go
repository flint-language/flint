package vm

import (
	"flint/internal/bytecode"
	"fmt"
)

type VM struct {
	chunk *bytecode.Chunk
	ip    int

	stack  Stack
	callStack []callFrame
	trace  bool
	globalFuncs []*bytecode.Function
}

func New(chunk *bytecode.Chunk) *VM {
	vm := &VM {
		chunk: chunk,
		ip: 0,
		stack: NewStack(),
		trace: true,
	}

	vm.globalFuncs = chunk.Funcs

	return vm
}

func (vm *VM) Run() {
	for {
		ip := vm.ip
		op := vm.readByte()

		if vm.trace {
			fmt.Printf("%04d %-6s -> %v %v\n", ip, op, vm.stack.data, vm.formatCallStack())
		}

		if vm.dispatch(op) {
			return
		}
	}
}

func (vm *VM) readByte() bytecode.OpCode {
	b := vm.chunk.Code[vm.ip]
	vm.ip++
	return b
}
