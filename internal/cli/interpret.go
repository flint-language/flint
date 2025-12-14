package cli

import (
	"flint/internal/bytecode"
	"flint/internal/vm"
)

func interpretFile(filename string) {
	prog, _ := loadAndParse(filename)
	chunk := bytecode.GenerateBytecode(prog)
	m := vm.New(chunk)
	m.Run()
}