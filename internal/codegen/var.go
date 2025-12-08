package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitVarDecl(b *ir.Block, e *parser.VarDeclExpr) value.Value {
	expr := cg.emitExpr(b, e.Value, false)
	alloc := b.NewAlloca(expr.Type())
	cg.locals[e.Name.Lexeme] = alloc
	b.NewStore(expr, alloc)
	return alloc
}
