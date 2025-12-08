package codegen

import (
	"flint/internal/lexer"
	"flint/internal/parser"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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

func (cg *CodeGen) emitList(b *ir.Block, e *parser.ListExpr) value.Value {
	exprs := make([]value.Value, 0)
	for _, elem := range e.Elements {
		expr := cg.emitExpr(b, elem, false)
		exprs = append(exprs, expr)
	}
	listType := exprs[0].Type()
	alloc := b.NewAlloca(types.NewArray(uint64(len(exprs)), listType))
	for idx, expr := range exprs {
		index := constant.NewInt(types.I32, int64(idx))
		elemPtr := b.NewGetElementPtr(listType, alloc, index)
		b.NewStore(expr, elemPtr)
	}
	return alloc
}

func (cg *CodeGen) emitIndex(b *ir.Block, e *parser.IndexExpr) value.Value {
	target := cg.emitExpr(b, e.Target, false)
	indexVal := cg.emitExpr(b, e.Index, false)
	switch e.Target.(type) {
	case *parser.TupleExpr, *parser.Identifier:
		idxConst, ok := indexVal.(*constant.Int)
		if !ok {
			panic("Index must be constant int for tuple")
		}
		ptr, ok := target.(*ir.InstAlloca)
		if !ok {
			ptrType := target.Type()
			ptr = b.NewAlloca(ptrType)
			b.NewStore(target, ptr)
		}
		structType := ptr.Type().(*types.PointerType).ElemType.(*types.StructType)
		elemPtr := b.NewGetElementPtr(
			structType,
			ptr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, idxConst.X.Int64()),
		)
		return b.NewLoad(structType.Fields[idxConst.X.Int64()], elemPtr)
	default:
		elemType := target.Type().(*types.PointerType).ElemType
		elemPtr := b.NewGetElementPtr(
			elemType,
			target,
			indexVal,
		)
		return b.NewLoad(elemType, elemPtr)
	}
}

func (cg *CodeGen) emitTuple(b *ir.Block, e *parser.TupleExpr) value.Value {
	exprs := make([]value.Value, len(e.Elements))
	tTypes := make([]types.Type, len(e.Elements))
	for i, elem := range e.Elements {
		exprs[i] = cg.emitExpr(b, elem, false)
		tTypes[i] = exprs[i].Type()
	}
	tupleType := types.NewStruct(tTypes...)
	alloc := b.NewAlloca(tupleType)
	for i, expr := range exprs {
		index := constant.NewInt(types.I32, int64(i))
		elemPtr := b.NewGetElementPtr(tupleType, alloc, constant.NewInt(types.I32, 0), index)
		b.NewStore(expr, elemPtr)
	}
	return alloc
}
