package codegen

import (
	"flint/internal/lexer"
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitInfix(b *ir.Block, e *parser.InfixExpr) value.Value {
	l := cg.emitExpr(b, e.Left, false)
	r := cg.emitExpr(b, e.Right, false)

	switch e.Operator.Kind {
	case lexer.Plus:
		return b.NewAdd(l, r)
	case lexer.Minus:
		return b.NewSub(l, r)
	case lexer.Star:
		return b.NewMul(l, r)
	case lexer.Slash:
		return b.NewSDiv(l, r)
	case lexer.Percent:
		return b.NewSRem(l, r)
	case lexer.Less:
		return b.NewICmp(enum.IPredSLT, l, r)
	case lexer.Greater:
		return b.NewICmp(enum.IPredSGT, l, r)
	case lexer.LessEqual:
		return b.NewICmp(enum.IPredSLE, l, r)
	case lexer.GreaterEqual:
		return b.NewICmp(enum.IPredSGE, l, r)
	case lexer.EqualEqual:
		return b.NewICmp(enum.IPredEQ, l, r)
	case lexer.NotEqual:
		return b.NewICmp(enum.IPredNE, l, r)
	case lexer.PlusDot:
		return b.NewFAdd(l, r)
	case lexer.MinusDot:
		return b.NewFSub(l, r)
	case lexer.StarDot:
		return b.NewFMul(l, r)
	case lexer.SlashDot:
		return b.NewFDiv(l, r)
	case lexer.LessDot:
		return b.NewFCmp(enum.FPredOLT, l, r)
	case lexer.GreaterDot:
		return b.NewFCmp(enum.FPredOGT, l, r)
	case lexer.LessEqualDot:
		return b.NewFCmp(enum.FPredOLE, l, r)
	case lexer.GreaterEqualDot:
		return b.NewFCmp(enum.FPredOGE, l, r)
	case lexer.LtGt:
		// return cg.emitConcat(b, e, l, r)
		return nil
	case lexer.AmperAmper:
		return b.NewAnd(l, r)
	case lexer.VbarVbar:
		return b.NewOr(l, r)
	}
	panic("unsupported operator")
}
