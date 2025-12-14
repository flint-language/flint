package bytecode

import (
	"flint/internal/parser"
)

type Generator struct {
	chunk       *Chunk
	globalChunk *Chunk
}

func GenerateBytecode(prog *parser.Program) *Chunk {
	chunk := &Chunk{}
	gen := &Generator{
		chunk:       chunk,
		globalChunk: chunk,
	}
	for _, expr := range prog.Exprs {
		if fn, ok := expr.(*parser.FuncDeclExpr); ok {
			gen.emitExpr(fn)
		}
	}
	for _, expr := range prog.Exprs {
		if _, ok := expr.(*parser.FuncDeclExpr); ok {
			continue
		}
		gen.emitExpr(expr)
	}
	mainIdx := -1
	for i, fn := range chunk.Funcs {
		if fn.Name == "main" {
			mainIdx = i
			break
		}
	}
	if mainIdx != -1 {
		gen.emit(OP_CALL)
		gen.emit(OpCode(mainIdx))
	}

	gen.emit(OP_HALT)
	return chunk
}
