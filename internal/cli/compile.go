package cli

import (
	"flint/internal/codegen"
	"fmt"
)

func compileFile(filename string) {
	fmt.Println("Compiling " + filename)
	prog, _ := loadAndParse(filename)
	ir := codegen.GenerateLLVM(prog)
	fmt.Println(ir)
}
