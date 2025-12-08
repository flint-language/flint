package codegen

import (
	"path/filepath"
)

func (cg *CodeGen) initModuleHeaders(sourceFile string) {
	cg.mod.SourceFilename = filepath.Base(sourceFile)
	cg.mod.TargetTriple = detectHostTriple()
}
