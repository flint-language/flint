package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"flint/internal/parser"
)

type CodeGen struct {
	mod      *ir.Module
	strIndex int
}

func GenerateLLVM(prog *parser.Program) string {
	cg := &CodeGen{
		mod: ir.NewModule(),
	}

	for _, e := range prog.Exprs {
		switch n := e.(type) {
		case *parser.FuncDeclExpr:
			cg.emitFunction(n)
		case *parser.IntLiteral, *parser.FloatLiteral, *parser.BoolLiteral,
			*parser.ByteLiteral, *parser.StringLiteral:
			cg.emitTopLiteral(n)
		default:
			panic("unsupported top-level expr")
		}
	}

	return cg.mod.String()
}

func (cg *CodeGen) emitFunction(fn *parser.FuncDeclExpr) {
	name := fn.Name.Lexeme
	ret := cg.resolveType(fn.Ret)
	if name == "main" {
		ret = types.I32
	}

	params := []*ir.Param{}
	for _, p := range fn.Params {
		params = append(params, ir.NewParam(p.Name.Lexeme, cg.resolveType(p.Type)))
	}

	mainFn := cg.mod.NewFunc(name, ret, params...)
	entry := mainFn.NewBlock("entry")

	if fn.Body == nil {
		cg.emitDefaultReturn(entry, ret, name == "main")
		return
	}

	block := fn.Body.(*parser.BlockExpr)
	val := cg.emitBlock(entry, block)

	if name == "main" {
		entry.NewRet(constant.NewInt(types.I32, 0))
		return
	}

	if val == nil {
		cg.emitDefaultReturn(entry, ret, false)
	} else {
		entry.NewRet(val)
	}
}

func (cg *CodeGen) emitBlock(b *ir.Block, blk *parser.BlockExpr) value.Value {
	var last value.Value

	for _, e := range blk.Exprs {
		last = cg.emitExpr(b, e)
	}

	return last
}
func (cg *CodeGen) emitExpr(_ *ir.Block, e parser.Expr) value.Value {
	switch v := e.(type) {
	case *parser.IntLiteral:
		return constant.NewInt(types.I64, v.Value)

	case *parser.FloatLiteral:
		return constant.NewFloat(types.Double, v.Value)

	case *parser.BoolLiteral:
		if v.Value {
			return constant.NewInt(types.I1, 1)
		}
		return constant.NewInt(types.I1, 0)

	case *parser.ByteLiteral:
		return constant.NewInt(types.I8, int64(v.Value))

	case *parser.StringLiteral:
		return cg.emitString(v)

	default:
		panic("unsupported expression type")
	}
}

func (cg *CodeGen) emitString(v *parser.StringLiteral) value.Value {
	label := cg.newStrLabel()

	str := constant.NewCharArrayFromString(v.Value + "\x00")

	global := cg.mod.NewGlobalDef(label, str)
	global.Immutable = true
	global.Align = 1

	zero := constant.NewInt(types.I32, 0)
	return constant.NewGetElementPtr(
		str.Typ,
		global,
		zero,
		zero,
	)
}

func (cg *CodeGen) newStrLabel() string {
	name := ".str." + string(rune('a'+cg.strIndex))
	cg.strIndex++
	return name
}

func (cg *CodeGen) resolveType(t parser.Expr) types.Type {
	if t == nil {
		return types.Void
	}

	ty := t.(*parser.TypeExpr)

	switch ty.Name {
	case "Int":
		return types.I64
	case "Float":
		return types.Double
	case "Bool":
		return types.I1
	case "Byte":
		return types.I8
	case "String":
		return types.I8Ptr
	case "Nil":
		return types.Void
	default:
		return types.I64
	}
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

func (cg *CodeGen) emitTopLiteral(e parser.Expr) {
	fn := cg.mod.NewFunc("main", types.I32)
	b := fn.NewBlock("entry")

	val := cg.emitExpr(b, e)

	if val.Type().Equal(types.I64) {
		val = b.NewTrunc(val, types.I32)
	}

	b.NewRet(val)
}
