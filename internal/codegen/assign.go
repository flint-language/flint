package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitAssign(b *ir.Block, e *parser.AssignExpr) value.Value {
	expr := cg.emitExpr(b, e.Value, false)
	alloc := cg.locals[e.Name.Name]
	b.NewStore(expr, alloc)
	return alloc
}
