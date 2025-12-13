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
	Name       string
	Kind       SymbolKind
	Type       string
	CurriedSig string
	Line       int
}

var symbols = map[string][]Symbol{}

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
				typ = node.Type.NodeType()
			}
			syms = append(syms, Symbol{
				Name: node.Name.Lexeme,
				Kind: VariableSymbol,
				Type: typ,
				Line: node.Name.Line,
			})
		case *parser.FuncDeclExpr:
			curriedTypes := make([]string, len(node.Params))
			for i, p := range node.Params {
				t := "_"
				if p.Type != nil {
					t = p.Type.NodeType()
				}
				curriedTypes[i] = t
			}
			retType := "_"
			if node.Ret != nil {
				retType = node.Ret.NodeType()
			}
			curriedSig := strings.Join(curriedTypes, " -> ") + " -> " + retType
			syms = append(syms, Symbol{
				Name:       node.Name.Lexeme,
				Kind:       FunctionSymbol,
				CurriedSig: curriedSig,
				Type:       node.Name.Lexeme,
				Line:       node.Name.Line,
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
