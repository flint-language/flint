package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitExpr(b *ir.Block, e parser.Expr, isTail bool) value.Value {
	switch v := e.(type) {
	case *parser.IntLiteral:
		return constant.NewInt(cg.platformIntType(), v.Value)
	case *parser.FloatLiteral:
		return constant.NewFloat(cg.platformFloatType(), v.Value)
	case *parser.BoolLiteral:
		if v.Value {
			return constant.NewInt(types.I1, 1)
		}
		return constant.NewInt(types.I1, 0)
	case *parser.ByteLiteral:
		return constant.NewInt(types.I8, int64(v.Value))
	case *parser.StringLiteral:
		return cg.emitString(v)
	case *parser.CallExpr:
		return cg.emitCall(b, v, isTail)
	case *parser.Identifier:
		ptr := cg.locals[v.Name]
		if ptr == nil {
			panic("undefined variable: " + v.Name)
		}
		return b.NewLoad(ptr.Type().(*types.PointerType).ElemType, ptr)
	case *parser.InfixExpr:
		return cg.emitInfix(b, v)
	case *parser.IfExpr:
		return cg.emitIf(b, v)
	case *parser.MatchExpr:
		return cg.emitMatch(b, v, isTail)
	case *parser.VarDeclExpr:
		return cg.emitVarDecl(b, v)
	case *parser.AssignExpr:
		return cg.emitAssign(b, v)
	case *parser.PrefixExpr:
		return cg.emitPrefix(b, v)
	default:
		panic("unsupported expression type")
	}
}
