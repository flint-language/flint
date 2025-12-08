package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitFunction(fn *parser.FuncDeclExpr) {
	name := fn.Name.Lexeme
	mainFn := cg.funcs[name]
	if len(fn.Decorators) != 0 && fn.Decorators[0].Name == "external" {
		cg.emitExternalFunction(fn, 0, name, mainFn)
		return
	}
	cg.locals = map[string]value.Value{}
	entry := mainFn.NewBlock("entry")
	for _, param := range mainFn.Params {
		alloc := entry.NewAlloca(param.Type())
		entry.NewStore(param, alloc)
		cg.locals[param.Name()] = alloc
	}
	if fn.Body == nil {
		if name == "main" {
			entry.NewRet(constant.NewInt(types.I32, 0))
		} else if mainFn.Sig.RetType.Equal(types.Void) {
			entry.NewRet(nil)
		} else {
			cg.emitDefaultReturn(entry, mainFn.Sig.RetType, false)
		}
		return
	}
	block := fn.Body.(*parser.BlockExpr)
	lastVal := cg.emitBlock(entry, block, fn.Recursion)
	retTy := mainFn.Sig.RetType
	if name == "main" {
		exit := mainFn.NewBlock("main.exit")
		for _, bb := range mainFn.Blocks {
			if bb.Term == nil && bb != exit {
				bb.NewBr(exit)
			}
		}
		exit.NewRet(constant.NewInt(types.I32, 0))
	} else if retTy.Equal(types.Void) {
		if entry.Term == nil {
			entry.NewRet(nil)
		}
	} else {
		if lastVal != nil {
			if b := parentBlockOfValue(lastVal); b != nil && b.Term == nil {
				b.NewRet(lastVal)
			} else if entry.Term == nil {
				entry.NewRet(lastVal)
			}
		} else if entry.Term == nil {
			cg.emitDefaultReturn(entry, retTy, false)
		}
	}
}

func (cg *CodeGen) emitBlock(b *ir.Block, blk *parser.BlockExpr, isTail bool) value.Value {
	var last value.Value
	n := len(blk.Exprs)
	for i, e := range blk.Exprs {
		tail := isTail && i == n-1
		last = cg.emitExpr(b, e, tail)
	}
	return last
}

func (cg *CodeGen) emitDefaultReturn(b *ir.Block, ret types.Type, isMain bool) {
	if isMain {
		b.NewRet(constant.NewInt(types.I32, 0))
		return
	}
	switch t := ret.(type) {
	case *types.IntType:
		b.NewRet(constant.NewInt(t, 0))
	case *types.FloatType:
		b.NewRet(constant.NewFloat(t, 0))
	case *types.PointerType:
		b.NewRet(constant.NewNull(t))
	case *types.VoidType:
		b.NewRet(nil)
	default:
		panic("unsupported return type")
	}
}
