package interpreter

import "flint/internal/parser"

func Eval(expr parser.Expr, env *Env) Value {
	switch e := expr.(type) {
	case *parser.IntLiteral:
		return Int(e.Value)
	case *parser.FloatLiteral:
		return Float(e.Value)
	case *parser.UnsignedLiteral:
		return Unsigned(e.Value)
	case *parser.StringLiteral:
		return String(e.Value)
	case *parser.BoolLiteral:
		return Bool(e.Value)
	case *parser.InfixExpr:
		left := Eval(e.Left, env)
		right := Eval(e.Right, env)
		return evalInfix(e.Operator.Lexeme, left, right)
	default:
		return Nil()
	}
}
