package codegen

import (
	"flint/internal/lexer"
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitPrefix(b *ir.Block, e *parser.PrefixExpr) value.Value {
	expr := cg.emitExpr(b, e.Right, false)
	ty := expr.Type()

	switch e.Operator.Kind {
	case lexer.Minus:
		intType := ty.(*types.IntType)
		return b.NewSub(constant.NewInt(intType, 0), expr)
	case lexer.MinusDot:
		floatType := ty.(*types.FloatType)
		return b.NewFSub(constant.NewFloat(floatType, 0), expr)
	case lexer.Bang:
		return b.NewXor(constant.True, expr)
	}
	return nil
}
