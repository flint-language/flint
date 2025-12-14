package vm

import (
	"flint/internal/bytecode"
	"flint/internal/runtime"
	"fmt"
)

func (vm *VM) intBinOp(fn func(a, b int64) int64) {
	b := vm.stack.Pop().AsInt()
	a := vm.stack.Pop().AsInt()
	vm.stack.Push(runtime.Int(fn(a, b)))
}

func (vm *VM) floatBinOp(fn func(a, b float64) float64) {
	b := vm.stack.Pop().AsFloat()
	a := vm.stack.Pop().AsFloat()
	vm.stack.Push(runtime.Float(fn(a, b)))
}

func (vm *VM) intCmp(fn func(a, b int64) bool) {
	b := vm.stack.Pop().AsInt()
	a := vm.stack.Pop().AsInt()
	vm.stack.Push(runtime.Bool(fn(a, b)))
}

func (vm *VM) floatCmp(fn func(a, b float64) bool) {
	b := vm.stack.Pop().AsFloat()
	a := vm.stack.Pop().AsFloat()
	vm.stack.Push(runtime.Bool(fn(a, b)))
}

func (vm *VM) dispatch(op bytecode.OpCode) bool {
	switch op {
	case bytecode.OP_CONST:
		idx := vm.readByte()
		vm.stack.Push(vm.chunk.Consts[idx])
	case bytecode.OP_ADD:
		vm.intBinOp(func(a, b int64) int64 { return a + b })
	case bytecode.OP_SUB:
		vm.intBinOp(func(a, b int64) int64 { return a - b })
	case bytecode.OP_MUL:
		vm.intBinOp(func(a, b int64) int64 { return a * b })
	case bytecode.OP_DIV:
		vm.intBinOp(func(a, b int64) int64 {
			if b == 0 {
				runtime.MathError("division by zero")
			}
			return a / b
		})
	case bytecode.OP_MOD:
		vm.intBinOp(func(a, b int64) int64 { return a % b })
	case bytecode.OP_LT:
		vm.intCmp(func(a, b int64) bool { return a < b })
	case bytecode.OP_GT:
		vm.intCmp(func(a, b int64) bool { return a > b })
	case bytecode.OP_LE:
		vm.intCmp(func(a, b int64) bool { return a <= b })
	case bytecode.OP_GE:
		vm.intCmp(func(a, b int64) bool { return a >= b })
	case bytecode.OP_FADD:
		vm.floatBinOp(func(a, b float64) float64 { return a + b })
	case bytecode.OP_FSUB:
		vm.floatBinOp(func(a, b float64) float64 { return a - b })
	case bytecode.OP_FMUL:
		vm.floatBinOp(func(a, b float64) float64 { return a * b })
	case bytecode.OP_FDIV:
		vm.floatBinOp(func(a, b float64) float64 { return a / b })
	case bytecode.OP_FLT:
		vm.floatCmp(func(a, b float64) bool { return a < b })
	case bytecode.OP_FGT:
		vm.floatCmp(func(a, b float64) bool { return a > b })
	case bytecode.OP_FLE:
		vm.floatCmp(func(a, b float64) bool { return a <= b })
	case bytecode.OP_FGE:
		vm.floatCmp(func(a, b float64) bool { return a >= b })
	case bytecode.OP_EQ:
		b := vm.stack.Pop()
		a := vm.stack.Pop()
		vm.stack.Push(runtime.Bool(a.Equals(b)))
	case bytecode.OP_NEQ:
		b := vm.stack.Pop()
		a := vm.stack.Pop()
		vm.stack.Push(runtime.Bool(!a.Equals(b)))
	case bytecode.OP_JUMP:
		vm.ip = int(vm.readU16())
	case bytecode.OP_JUMP_IF_FALSE:
		target := vm.readU16()
		cond := vm.stack.Pop()
		if !cond.IsTruthy() {
			vm.ip = int(target)
		}
	case bytecode.OP_HALT:
		return true
	case bytecode.OP_CALL:
		fnIdx := int(vm.readByte())
		fn := vm.globalFuncs[fnIdx]
		args := make([]runtime.Value, fn.Params)
		for i := fn.Params - 1; i >= 0; i-- {
			args[i] = vm.stack.Pop()
		}
		vm.callStack = append(vm.callStack, callFrame{
			ip:         vm.ip,
			chunk:      vm.chunk,
			funcIndex: fnIdx,
			stackStart: len(vm.stack.data),
		})
		vm.chunk = fn.Chunk
		vm.ip = 0
		for _, arg := range args {
			vm.stack.Push(arg)
		}
	case bytecode.OP_RETURN:
		ret := vm.stack.Pop()
		frame := vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		vm.stack.data = vm.stack.data[:frame.stackStart]
		vm.stack.Push(ret)
		vm.chunk = frame.chunk
		vm.ip = frame.ip
	default:
		runtime.RuntimeError(fmt.Sprintf("unknown opcode %d", op))
	}
	return false
}
