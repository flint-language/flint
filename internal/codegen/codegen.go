package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"flint/internal/parser"
)

type CodeGen struct {
	mod              *ir.Module
	strIndex         int
	strGlobals       map[string]*ir.Global
	globalMatchCount int

	locals map[string]value.Value
	funcs  map[string]*ir.Func
}

func GenerateLLVM(prog *parser.Program, sourceFile string) string {
	cg := &CodeGen{
		mod:        ir.NewModule(),
		locals:     map[string]value.Value{},
		funcs:      map[string]*ir.Func{},
		strGlobals: map[string]*ir.Global{},
	}
	cg.initModuleHeaders(sourceFile)
	for _, e := range prog.Exprs {
		if fn, ok := e.(*parser.FuncDeclExpr); ok {
			name := fn.Name.Lexeme
			ret := cg.resolveType(fn.Ret)
			if name == "main" {
				ret = types.I32
			}
			params := []*ir.Param{}
			for _, p := range fn.Params {
				params = append(params, ir.NewParam(p.Name.Lexeme, cg.resolveType(p.Type)))
			}
			cg.funcs[name] = cg.mod.NewFunc(name, ret, params...)
		}
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
