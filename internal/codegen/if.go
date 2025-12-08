package codegen

import (
	"flint/internal/parser"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitIf(b *ir.Block, i *parser.IfExpr, isTail bool) value.Value {
	cond := cg.emitExpr(b, i.Cond, false)
	parent := b.Parent
	thenBlock := parent.NewBlock("if.then")
	elseBlock := parent.NewBlock("if.else")
	mergeBlock := parent.NewBlock("if.merge")
	b.NewCondBr(cond, thenBlock, elseBlock)
	thenVal := cg.emitIfBody(thenBlock, i.Then, false)
	if thenBlock.Term == nil {
		thenBlock.NewBr(mergeBlock)
	}
	elseVal := cg.emitIfBody(elseBlock, i.Else, false)
	if elseBlock.Term == nil {
		elseBlock.NewBr(mergeBlock)
	}
	var phiType types.Type
	if thenVal != nil {
		phiType = thenVal.Type()
	}
	if elseVal != nil {
		if phiType == nil {
			phiType = elseVal.Type()
		} else if !elseVal.Type().Equal(phiType) {
			panic(fmt.Sprintf("if branch type mismatch: %v vs %v", elseVal.Type(), phiType))
		}
	}
	if phiType == nil {
		mergeBlock.NewRet(nil)
		return nil
	}
	if phiType.Equal(types.Void) {
		return mergeBlock
	}
	var incomings []*ir.Incoming
	if thenVal != nil {
		incomings = append(incomings, &ir.Incoming{X: thenVal, Pred: thenBlock})
	} else {
		incomings = append(incomings, &ir.Incoming{X: constant.NewUndef(phiType), Pred: thenBlock})
	}
	if elseVal != nil {
		incomings = append(incomings, &ir.Incoming{X: elseVal, Pred: elseBlock})
	} else {
		incomings = append(incomings, &ir.Incoming{X: constant.NewUndef(phiType), Pred: elseBlock})
	}
	if referencesBlock(b, mergeBlock) {
		incomings = append(incomings, ir.NewIncoming(constant.NewUndef(phiType), b))
	}
	phi := mergeBlock.NewPhi(incomings...)
	if mergeBlock.Term == nil {
		if phiType.Equal(types.Void) {
			mergeBlock.NewRet(nil)
		} else {
			mergeBlock.NewRet(phi)
		}
	} else if isTail {
		mergeBlock.NewRet(phi)
	}
	return phi
}

func (cg *CodeGen) emitIfBody(b *ir.Block, body parser.Expr, isTail bool) value.Value {
	switch x := body.(type) {
	case *parser.BlockExpr:
		return cg.emitBlock(b, x, isTail)
	default:
		return cg.emitExpr(b, body, isTail)
	}
}
