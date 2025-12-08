package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGen) emitTopLiteral(e parser.Expr) {
	fn := cg.mod.NewFunc("main", types.I32)
	b := fn.NewBlock("entry")
	_ = cg.emitExpr(b, e, true)
	b.NewRet(constant.NewInt(types.I32, 0))
}
