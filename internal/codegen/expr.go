package codegen

import (
	"flint/internal/parser"
	"fmt"

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
		return cg.emitIf(b, v, isTail)
	case *parser.MatchExpr:
		return cg.emitMatch(b, v, isTail)
	case *parser.VarDeclExpr:
		return cg.emitVarDecl(b, v)
	case *parser.AssignExpr:
		return cg.emitAssign(b, v)
	case *parser.PrefixExpr:
		return cg.emitPrefix(b, v)
	case *parser.ListExpr:
		return cg.emitList(b, v)
	case *parser.TupleExpr:
		return cg.emitTuple(b, v)
	case *parser.IndexExpr:
		return cg.emitIndex(b, v)
	case *parser.FuncDeclExpr:
		unique := fmt.Sprintf("%s$%d", v.Name.Lexeme, len(cg.funcs))
		retTy := cg.resolveType(v.Ret)
		params := []*ir.Param{}
		for _, p := range v.Params {
			params = append(params, ir.NewParam(p.Name.Lexeme, cg.resolveType(p.Type)))
		}
		irfn := cg.mod.NewFunc(unique, retTy, params...)
		cg.funcs[v.Name.Lexeme] = irfn
		savedLocals := cg.locals
		cg.locals = map[string]value.Value{}
		entry := irfn.NewBlock("entry")
		for _, param := range irfn.Params {
			alloc := entry.NewAlloca(param.Type())
			entry.NewStore(param, alloc)
			cg.locals[param.Name()] = alloc
		}
		block := v.Body.(*parser.BlockExpr)
		lastVal := cg.emitBlock(entry, block, v.Recursion)
		if retTy.Equal(types.Void) {
			entry.NewRet(nil)
		} else if lastVal != nil {
			if parent := parentBlockOfValue(lastVal); parent != nil && parent.Term == nil {
				parent.NewRet(lastVal)
			} else if entry.Term == nil {
				entry.NewRet(lastVal)
			}
		} else {
			entry.NewRet(constant.NewInt(types.I64, 0))
		}
		cg.locals = savedLocals
		return nil
	default:
		panic("unsupported expression type")
	}
}
