package cli

import (
	"flint/internal/codegen"
	"fmt"
	"os"
	"strings"
)

func compileFile(filename string) {
	prog, _ := loadAndParse(filename)
	ir := codegen.GenerateLLVM(prog)
	outFile := filename
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		outFile = filename[:idx]
	}
	outFile += ".ll"
	if err := os.WriteFile(outFile, []byte(ir), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing LLVM IR:", err)
		return
	}
}
