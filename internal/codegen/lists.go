package codegen

import (
	"flint/internal/parser"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

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
	index := cg.emitExpr(b, e.Index, false)

	elemType := (target.Type().(*types.PointerType)).ElemType
	fmt.Println(target, index, elemType)

	elemPtr := b.NewGetElementPtr(elemType, target, index)

	return b.NewLoad(elemType, elemPtr)
}

func (cg *CodeGen) emitTuple(b *ir.Block, e *parser.TupleExpr) value.Value {
	exprs := make([]value.Value, 0)
	tTypes := make([]types.Type, 0)

	for _, elem := range e.Elements {
		expr := cg.emitExpr(b, elem, false)
		exprs = append(exprs, expr)
		tTypes = append(tTypes, expr.Type())
	}
	tupleType := types.NewStruct(tTypes...)
	alloc := b.NewAlloca(tupleType)

	for idx, expr := range exprs {
		index := constant.NewInt(types.I32, int64(idx))
		elemPtr := b.NewGetElementPtr(tupleType, alloc, constant.NewInt(types.I32, 0), index)
		b.NewStore(expr, elemPtr)
	}
	return alloc
}
