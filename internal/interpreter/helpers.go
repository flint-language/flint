package interpreter

import (
	"flint/internal/lexer"
	"flint/internal/parser"
	"fmt"
)

func evalInfix(op string, left, right Value) Value {
	switch op {
	case "+":
		return Int(left.Int + right.Int)
	case "-":
		return Int(left.Int - right.Int)
	case "*":
		return Int(left.Int * right.Int)
	case "/":
		return Int(left.Int / right.Int)
	case "+.":
		return Float(left.Float + right.Float)
	case "-.":
		return Float(left.Float - right.Float)
	case "*.":
		return Float(left.Float * right.Float)
	case "/.":
		return Float(left.Float / right.Float)
	default:
		panic("unknown infix operator: " + op)
	}
}

func RunReplLine(input string, env *Env) {
	tokens, err := lexer.Tokenize(input, "<repl>")
	if err != nil {
		fmt.Println("lex error:", err)
		return
	}
	p, _ := parser.ParseProgram(tokens)
	value := Eval(p.Exprs[0], env)
	fmt.Println(PrintReplResult(value))
}
