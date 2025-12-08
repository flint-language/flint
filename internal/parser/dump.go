package parser

import (
	"fmt"
	"strings"
)

func DumpExpr(e Expr) string {
	return dump(e, "", true)
}

func node(indent string, last bool, label string) (string, string) {
	branch := "├─ "
	next := indent + "│  "
	if last {
		branch = "└─ "
		next = indent + "   "
	}
	return indent + branch + label + "\n", next
}

func dump(e Expr, indent string, last bool) string {
	switch n := e.(type) {
	case *Identifier:
		line, _ := node(indent, last, "Identifier "+n.Name)
		return line
	case *IntLiteral:
		line, _ := node(indent, last, fmt.Sprintf("Int %d", n.Value))
		return line
	case *FloatLiteral:
		line, _ := node(indent, last, fmt.Sprintf("Float %g", n.Value))
		return line
	case *StringLiteral:
		line, _ := node(indent, last, fmt.Sprintf("String %q", n.Value))
		return line
	case *ByteLiteral:
		line, _ := node(indent, last, fmt.Sprintf("Byte '%c'", n.Value))
		return line
	case *BoolLiteral:
		line, _ := node(indent, last, fmt.Sprintf("Bool %t", n.Value))
		return line
	case *PrefixExpr:
		line, next := node(indent, last, "Prefix "+n.Operator.Lexeme)
		return line + dump(n.Right, next, true)
	case *InfixExpr:
		line, next := node(indent, last, "Infix "+n.Operator.Lexeme)
		return line +
			dump(n.Left, next, false) +
			dump(n.Right, next, true)
	case *CallExpr:
		line, next := node(indent, last, "Call")
		var out strings.Builder
		out.WriteString(line)
		cLine, cNext := node(next, false, "Callee")
		out.WriteString(cLine)
		out.WriteString(dump(n.Callee, cNext, true))
		if len(n.Args) > 0 {
			aLine, aNext := node(next, true, "Args")
			out.WriteString("\n")
			out.WriteString(aLine)
			for i, arg := range n.Args {
				out.WriteString(dump(arg, aNext, i == len(n.Args)-1))
			}
		}
		return out.String()
	case *PipelineExpr:
		line, next := node(indent, last, "Pipeline")
		return line +
			dump(n.Left, next, false) +
			dump(n.Right, next, true)
	case *QualifiedExpr:
		line, next := node(indent, last, "Qualified")
		return line +
			dump(n.Left, next, false) +
			nodeWith(next, true, "Identifier "+n.Right.Lexeme)
	case *FieldAccessExpr:
		line, next := node(indent, last, "FieldAccess")
		return line +
			dump(n.Left, next, false) +
			nodeWith(next, true, "Identifier "+n.Right)
	case *TupleExpr:
		line, next := node(indent, last, "Tuple")
		var out strings.Builder
		out.WriteString(line)
		for i, e := range n.Elements {
			out.WriteString(dump(e, next, i == len(n.Elements)-1))
		}
		return out.String()
	case *ListExpr:
		line, next := node(indent, last, "List")
		var out strings.Builder
		out.WriteString(line)
		for i, e := range n.Elements {
			out.WriteString(dump(e, next, i == len(n.Elements)-1))
		}
		return out.String()
	case *BlockExpr:
		line, next := node(indent, last, "Block")
		var out strings.Builder
		out.WriteString(line)
		for i, e := range n.Exprs {
			out.WriteString(dump(e, next, i == len(n.Exprs)-1))
		}
		return out.String()
	case *IfExpr:
		line, next := node(indent, last, "IfExpr")
		var out strings.Builder
		out.WriteString(line)
		cLine, cNext := node(next, false, "Cond")
		out.WriteString(cLine)
		out.WriteString(dump(n.Cond, cNext, true))
		tLine, tNext := node(next, n.Else == nil, "Then")
		out.WriteString("\n")
		out.WriteString(tLine)
		out.WriteString(dump(n.Then, tNext, true))
		if n.Else != nil {
			eLine, eNext := node(next, true, "Else")
			out.WriteString("\n")
			out.WriteString(eLine)
			out.WriteString(dump(n.Else, eNext, true))
		}
		return out.String()
	case *MatchExpr:
		line, next := node(indent, last, "MatchExpr")
		var out strings.Builder
		out.WriteString(line)
		vLine, vNext := node(next, false, "Value")
		out.WriteString(vLine)
		out.WriteString(dump(n.Value, vNext, true))
		aLine, aNext := node(next, true, "Arms")
		out.WriteString("\n")
		out.WriteString(aLine)
		for i, arm := range n.Arms {
			armLast := i == len(n.Arms)-1
			armLine, armNext := node(aNext, armLast, "Arm")
			out.WriteString(armLine)
			pLine, pNext := node(armNext, false, "Pattern")
			out.WriteString(pLine)
			out.WriteString(dump(arm.Pattern, pNext, true))
			if arm.Guard != nil {
				gLine, gNext := node(armNext, false, "Guard")
				out.WriteString("\n")
				out.WriteString(gLine)
				out.WriteString(dump(arm.Guard, gNext, true))
			}
			bLine, bNext := node(armNext, true, "Body")
			out.WriteString("\n")
			out.WriteString(bLine)
			out.WriteString(dump(arm.Body, bNext, true))
			out.WriteString("\n")
		}
		return strings.TrimRight(out.String(), "\n")
	case *VarDeclExpr:
		line, next := node(indent, last, fmt.Sprintf("VarDecl name=%s mutable=%t", n.Name.Lexeme, n.Mutable))
		var out strings.Builder
		out.WriteString(line)
		if n.Type != nil {
			tLine, tNext := node(next, false, "Type")
			out.WriteString(tLine)
			out.WriteString(dump(n.Type, tNext, true))
		}
		vLine, vNext := node(next, true, "Value")
		out.WriteString("\n")
		out.WriteString(vLine)
		out.WriteString(dump(n.Value, vNext, true))
		return out.String()
	case *FuncDeclExpr:
		line, next := node(indent, last,
			fmt.Sprintf("FuncDecl name=%s pub=%t rec=%t",
				n.Name.Lexeme, n.Pub, n.Recursion),
		)
		var out strings.Builder
		out.WriteString(line)
		if len(n.Decorators) > 0 {
			dLine, dNext := node(next, true, "Decorators")
			out.WriteString("\n")
			out.WriteString(dLine)
			for i, dec := range n.Decorators {
				out.WriteString(dump(&dec, dNext, i == len(n.Decorators)-1))
			}
		}
		pLine, pNext := node(next, false, "Params")
		out.WriteString(pLine)
		for i, p := range n.Params {
			paramLine, pIndent := node(pNext, i == len(n.Params)-1, "Param name="+p.Name.Lexeme)
			out.WriteString(paramLine)
			if p.Type != nil {
				out.WriteString(dump(p.Type, pIndent, true))
			}
		}
		if n.Ret != nil {
			rLine, rNext := node(next, false, "ReturnType")
			out.WriteString("\n")
			out.WriteString(rLine)
			out.WriteString(dump(n.Ret, rNext, true))
		}
		bLine, bNext := node(next, true, "Body")
		out.WriteString("\n")
		out.WriteString(bLine)
		out.WriteString(dump(n.Body, bNext, true))
		return out.String()
	case *TypeDeclExpr:
		line, next := node(indent, last, "TypeDecl name="+n.Name.Lexeme)
		var out strings.Builder
		out.WriteString(line)
		if n.Body != nil {
			bLine, bNext := node(next, true, "Body")
			out.WriteString(bLine)
			out.WriteString(dump(n.Body, bNext, true))
		}
		return out.String()
	case *TypeExpr:
		line, next := node(indent, last, "Type "+n.Name)
		if n.Generic != nil {
			gLine, gNext := node(next, true, "Generic")
			return line + gLine + dump(n.Generic, gNext, true)
		}
		return line
	case *TupleTypeExpr:
		line, next := node(indent, last, "TupleType")
		var out strings.Builder
		out.WriteString(line)
		for i, t := range n.Types {
			out.WriteString(dump(t, next, i == len(n.Types)-1))
		}
		return out.String()
	case *RecordTypeExpr:
		line, next := node(indent, last, "RecordType "+n.Name.Lexeme)
		var out strings.Builder
		out.WriteString(line)
		for i, f := range n.Fields {
			fieldLine, fNext := node(next, i == len(n.Fields)-1, "Field "+f.Name.Lexeme)
			out.WriteString(fieldLine)
			if f.Type != nil {
				out.WriteString(dump(f.Type, fNext, true))
			}
		}
		return out.String()
	case *UseExpr:
		line, next := node(indent, last, "Use")
		var out strings.Builder
		out.WriteString(line)
		pathLine, _ := node(next, true, fmt.Sprintf("Path %v", n.Path))
		out.WriteString(pathLine)
		if len(n.Members) > 0 {
			mLine, _ := node(next, true, fmt.Sprintf("Members %v", n.Members))
			out.WriteString(mLine)
		}
		if n.Alias != "" {
			aLine, _ := node(next, true, "Alias "+n.Alias)
			out.WriteString(aLine)
		}
		return out.String()
	case *Decorator:
		line, next := node(indent, last, "Decorator "+n.Name)
		var out strings.Builder
		out.WriteString(line)
		if len(n.Args) > 0 {
			argsLine, argsNext := node(next, true, "Args")
			out.WriteString("\n")
			out.WriteString(argsLine)
			for i, arg := range n.Args {
				out.WriteString(dump(arg, argsNext, i == len(n.Args)-1))
			}
		}
		return out.String()
	case *AssignExpr:
		line, next := node(indent, last, "Assign")
		var out strings.Builder
		out.WriteString(line)
		lhsLine, lhsNext := node(next, false, "LHS")
		out.WriteString(lhsLine)
		out.WriteString(dump(n.Name, lhsNext, true))
		rhsLine, rhsNext := node(next, true, "RHS")
		out.WriteString("\n")
		out.WriteString(rhsLine)
		out.WriteString(dump(n.Value, rhsNext, true))
		return out.String()
	default:
		line, _ := node(indent, last, fmt.Sprintf("<unknown %T>", n))
		return line
	}
}

func nodeWith(indent string, last bool, label string) string {
	line, _ := node(indent, last, label)
	return line
}
