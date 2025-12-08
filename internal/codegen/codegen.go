package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"flint/internal/lexer"
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

func (cg *CodeGen) emitIf(b *ir.Block, i *parser.IfExpr) value.Value {
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
	thenVoid := thenVal == nil || thenVal.Type().Equal(types.Void)
	elseVoid := elseVal == nil || elseVal.Type().Equal(types.Void)
	if thenVoid && elseVoid {
		mergeBlock.NewBr(parent.NewBlock("if.continue"))
		return nil
	}
	phi := mergeBlock.NewPhi(
		&ir.Incoming{X: thenVal, Pred: thenBlock},
		&ir.Incoming{X: elseVal, Pred: elseBlock},
	)
	cont := parent.NewBlock("if.continue")
	mergeBlock.NewBr(cont)

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

func (cg *CodeGen) emitAssign(b *ir.Block, e *parser.AssignExpr) value.Value {
	expr := cg.emitExpr(b, e.Value, false)
	alloc := cg.locals[e.Name.Name]
	b.NewStore(expr, alloc)
	return alloc
}

func (cg *CodeGen) emitVarDecl(b *ir.Block, e *parser.VarDeclExpr) value.Value {
	expr := cg.emitExpr(b, e.Value, false)
	alloc := b.NewAlloca(expr.Type())
	cg.locals[e.Name.Lexeme] = alloc
	b.NewStore(expr, alloc)
	return alloc
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

func (cg *CodeGen) emitMatchBody(b *ir.Block, body parser.Expr, isTail bool) value.Value {
	switch x := body.(type) {
	case *parser.BlockExpr:
		return cg.emitBlock(b, x, isTail)
	case *parser.MatchExpr:
		return cg.emitMatch(b, x, isTail)
	default:
		return cg.emitExpr(b, body, isTail)
	}
}

func (cg *CodeGen) emitMatchCond(b *ir.Block, scr value.Value, pat parser.Expr, guard parser.Expr) value.Value {
	var baseCond value.Value
	switch p := pat.(type) {
	case *parser.IntLiteral:
		if _, ok := scr.Type().(*types.IntType); !ok {
			panic("match scrutinee is not an integer for integer pattern")
		}
		lit := constant.NewInt(types.I64, p.Value)
		baseCond = b.NewICmp(enum.IPredEQ, scr, lit)
	case *parser.BoolLiteral:
		if _, ok := scr.Type().(*types.IntType); !ok {
			panic("match scrutinee is not an integer for bool pattern")
		}
		var intVal int64
		if p.Value {
			intVal = 1
		} else {
			intVal = 0
		}
		lit := constant.NewInt(types.I1, intVal)
		baseCond = b.NewICmp(enum.IPredEQ, scr, lit)
	case *parser.Identifier:
		if p.Pos.Kind == lexer.Underscore {
			baseCond = constant.True
		} else {
			alloc := b.NewAlloca(scr.Type())
			b.NewStore(scr, alloc)
			cg.locals[p.Name] = alloc
			baseCond = constant.True
		}
	default:
		panic("unsupported match pattern: " + pat.NodeType())
	}
	if guard != nil {
		guardVal := cg.emitExpr(b, guard, false)
		if _, ok := guardVal.Type().(*types.IntType); ok {
			if guardVal.Type() != types.I1 {
				guardVal = b.NewTrunc(guardVal, types.I1)
			}
		} else {
			panic("guard expression is not a boolean")
		}
		baseCond = b.NewAnd(baseCond, guardVal)
	}
	return baseCond
}

func (cg *CodeGen) emitTopLiteral(e parser.Expr) {
	fn := cg.mod.NewFunc("main", types.I32)
	b := fn.NewBlock("entry")
	_ = cg.emitExpr(b, e, true)
	b.NewRet(constant.NewInt(types.I32, 0))
}
