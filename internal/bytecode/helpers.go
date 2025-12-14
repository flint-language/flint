package bytecode

import (
	"flint/internal/runtime"
)

func (g *Generator) emit(op OpCode) {
	g.chunk.Code = append(g.chunk.Code, op)
}

func (g *Generator) emitConst(v runtime.Value) {
	idx := len(g.chunk.Consts)
	g.chunk.Consts = append(g.chunk.Consts, v)
	g.emit(OP_CONST)
	g.emit(OpCode(idx))
}
