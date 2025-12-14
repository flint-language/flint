package bytecode

import (
	"flint/internal/parser"
	"flint/internal/runtime"
	"fmt"
)

func (g *Generator) emitExpr(e parser.Expr) {
	switch n := e.(type) {
	case *parser.IntLiteral:
		g.emitConst(runtime.Int(n.Value))
	case *parser.FloatLiteral:
		g.emitConst(runtime.Float(n.Value))
	case *parser.BoolLiteral:
		g.emitConst(runtime.Bool(n.Value))
	case *parser.StringLiteral:
		g.emitConst(runtime.String(n.Value))
	case *parser.InfixExpr:
		g.emitExpr(n.Left)
		g.emitExpr(n.Right)
		op := n.Operator.Lexeme
		switch op {
		case "+":
			g.emit(OP_ADD)
		case "-":
			g.emit(OP_SUB)
		case "*":
			g.emit(OP_MUL)
		case "/":
			g.emit(OP_DIV)
		case "%":
			g.emit(OP_MOD)
		case "+.":
			g.emit(OP_FADD)
		case "-.":
			g.emit(OP_FSUB)
		case "*.":
			g.emit(OP_FMUL)
		case "/.":
			g.emit(OP_FDIV)
		case "<":
			g.emit(OP_LT)
		case ">":
			g.emit(OP_GT)
		case "<=":
			g.emit(OP_LE)
		case ">=":
			g.emit(OP_GE)
		case "<.":
			g.emit(OP_FLT)
		case ">.":
			g.emit(OP_FGT)
		case "<=.":
			g.emit(OP_FLE)
		case ">=.":
			g.emit(OP_FGE)
		case "==":
			g.emit(OP_EQ)
		case "!=":
			g.emit(OP_NEQ)
		default:
			panic(fmt.Sprintf("unsupported infix operator %q", op))
		}

	case *parser.BlockExpr:
		for i, expr := range n.Exprs {
			isLast := i == len(n.Exprs)-1
			g.emitExpr(expr)
			if !isLast {
				g.emit(OP_POP)
			}
		}
	case *parser.FuncDeclExpr:
		funcChunk := &Chunk{}
		funcGen := &Generator{
			chunk:       funcChunk,
			globalChunk: g.globalChunk,
		}
		funcGen.emitExpr(n.Body)
		funcGen.emit(OP_RETURN)
		fn := &Function{
			Name:   n.Name.Lexeme,
			Chunk:  funcChunk,
			Params: len(n.Params),
		}
		g.chunk.Funcs = append(g.chunk.Funcs, fn)
	case *parser.CallExpr:
		for _, arg := range n.Args {
			g.emitExpr(arg)
		}
		fnIdx := -1
		callee := n.Callee.(*parser.Identifier).Name
		for i, fn := range g.globalChunk.Funcs {
			if fn.Name == callee {
				fnIdx = i
				break
			}
		}
		if fnIdx == -1 {
			panic("call to undefined function " + callee)
		}
		g.emit(OP_CALL)
		g.emit(OpCode(fnIdx))
	}
}
