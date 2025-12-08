package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitCall(b *ir.Block, c *parser.CallExpr, isTail bool) value.Value {
	id, ok := c.Callee.(*parser.Identifier)
	if !ok {
		panic("only simple function calls supported")
	}
	fn := cg.funcs[id.Name]
	if fn == nil {
		panic("undefined function: " + id.Name)
	}
	var args []value.Value
	for _, arg := range c.Args {
		args = append(args, cg.emitExpr(b, arg, false))
	}
	callInst := b.NewCall(fn, args...)
	if isTail {
		callerRet := b.Parent.Sig.RetType
		calleeRet := fn.Sig.RetType
		if calleeRet != nil && callerRet != nil && calleeRet.Equal(callerRet) {
			callInst.Tail = enum.TailTail
		}
	}
	return callInst
}
