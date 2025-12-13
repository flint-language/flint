package lsp

import (
	"flint/internal/parser"
	"strings"
)

type SymbolKind int

const (
	FunctionSymbol SymbolKind = iota + 1
	VariableSymbol
)

type Symbol struct {
	Name string
	Kind SymbolKind
	Type string
	Line int
}

func updateSymbols(uri string, prog *parser.Program) {
	syms := symbols[uri]
	syms = syms[:0]
	stack := make([]parser.Expr, len(prog.Exprs))
	copy(stack, prog.Exprs)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == nil {
			continue
		}
		switch node := n.(type) {
		case *parser.VarDeclExpr:
			typ := "_"
			if node.Type != nil {
				typ = formatType(node.Type)
			}
			syms = append(syms, Symbol{
				Name: node.Name.Lexeme,
				Kind: VariableSymbol,
				Type: typ,
				Line: node.Name.Line,
			})
		case *parser.FuncDeclExpr:
			paramTypes := make([]string, len(node.Params))
			for i, p := range node.Params {
				paramTypes[i] = formatType(p.Type)
			}
			retType := formatType(node.Ret)
			var sig strings.Builder
			sig.WriteString("pub fun " + node.Name.Lexeme + "(")
			for i, p := range node.Params {
				if i > 0 {
					sig.WriteString(", ")
				}
				sig.WriteString(p.Name.Lexeme + ": " + formatType(p.Type))
			}
			sig.WriteString(") " + retType)
			syms = append(syms, Symbol{
				Name: node.Name.Lexeme,
				Kind: FunctionSymbol,
				Type: sig.String(),
				Line: node.Name.Line,
			})
			if node.Body != nil {
				stack = append(stack, node.Body)
			}
		case *parser.BlockExpr:
			for i := len(node.Exprs) - 1; i >= 0; i-- {
				stack = append(stack, node.Exprs[i])
			}
		case *parser.AssignExpr:
			syms = append(syms, Symbol{
				Name: node.Name.Name,
				Kind: VariableSymbol,
				Type: node.Value.NodeType(),
			})
		}
	}
	symbols[uri] = syms
}
