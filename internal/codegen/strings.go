package codegen

import (
	"flint/internal/parser"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitString(v *parser.StringLiteral) value.Value {
	if g, ok := cg.strGlobals[v.Value]; ok {
		zero := constant.NewInt(types.I32, 0)
		return constant.NewGetElementPtr(g.Init.Type(), g, zero, zero)
	}
	label := cg.newStrLabel()
	str := constant.NewCharArrayFromString(v.Value + "\x00")
	global := cg.mod.NewGlobalDef(label, str)
	global.Immutable = true
	global.Align = 1
	cg.strGlobals[v.Value] = global
	zero := constant.NewInt(types.I32, 0)
	return constant.NewGetElementPtr(str.Typ, global, zero, zero)
}
