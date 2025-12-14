package vm

func (vm *VM) readU16() uint16 {
	hi := uint16(vm.readByte())
	lo := uint16(vm.readByte())
	return (hi << 8) | lo
}

func (vm *VM) formatCallStack() string {
	if len(vm.callStack) == 0 {
		return "[]"
	}

	out := "["
	for i, f := range vm.callStack {
		if i > 0 {
			out += " -> "
		}

		if f.funcIndex < 0 || f.funcIndex >= len(vm.globalFuncs) {
			out += "<root>"
		} else {
			out += vm.globalFuncs[f.funcIndex].Name
		}
	}
	out += "]"
	return out
}