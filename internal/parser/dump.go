package parser

import (
	"fmt"
	"strings"
)

func DumpExpr(e Expr) string {
	return dump(e, "", true)
}

func dump(e Expr, indent string, isLast bool) string {
	connector := "├─ "
	nextIndent := indent + "│  "
	if isLast {
		connector = "└─ "
		nextIndent = indent + "   "
	}
	switch n := e.(type) {
	case *Identifier:
		return fmt.Sprintf("%s%sIdentifier(%s)", indent, connector, n.Name)
	case *IntLiteral:
		return fmt.Sprintf("%s%sInt(%d)", indent, connector, n.Value)
	case *FloatLiteral:
		return fmt.Sprintf("%s%sFloat(%g)", indent, connector, n.Value)
	case *StringLiteral:
		return fmt.Sprintf("%s%sString(%s)", indent, connector, n.Value)
	case *ByteLiteral:
		return fmt.Sprintf("%s%sByte(%c)", indent, connector, n.Value)
	case *BoolLiteral:
		return fmt.Sprintf("%s%sBool(%t)", indent, connector, n.Value)
	case *PrefixExpr:
		return fmt.Sprintf("%s%sPrefix(%s)\n%s",
			indent, connector, n.Operator.Lexeme,
			dump(n.Right, nextIndent, true),
		)
	case *InfixExpr:
		return fmt.Sprintf("%s%sInfix(%s)\n%s\n%s",
			indent, connector, n.Operator.Lexeme,
			dump(n.Left, nextIndent, false),
			dump(n.Right, nextIndent, true),
		)
	case *CallExpr:
		out := fmt.Sprintf("%s%sCall\n%sCallee:\n%s\n%sArgs:\n",
			indent, connector, nextIndent, dump(n.Callee, nextIndent+"  ", true), nextIndent)
		for i, a := range n.Args {
			out += dump(a, nextIndent+"  ", i == len(n.Args)-1) + "\n"
		}
		return strings.TrimRight(out, "\n")
	case *ValDeclExpr:
		typeStr := ""
		if n.Type != nil {
			typeStr = ": " + dump(n.Type, "", true)
		}
		return fmt.Sprintf("%s%sValDecl(%s%s)\n%s",
			indent, connector, n.Name.Lexeme, typeStr,
			dump(n.Value, indent+"   ", true),
		)
	case *MutDeclExpr:
		typeStr := ""
		if n.Type != nil {
			typeStr = ": " + dump(n.Type, "", true)
		}
		return fmt.Sprintf("%s%sMutDecl(%s%s)\n%s",
			indent, connector, n.Name.Lexeme, typeStr,
			dump(n.Value, indent+"   ", true),
		)
	case *FuncDeclExpr:
		out := fmt.Sprintf("%s%sFuncDecl(pub=%t, rec=%t, name=%s)\n",
			indent, connector, n.Pub, n.Recursion, n.Name.Lexeme)
		out += fmt.Sprintf("%sParams:\n", nextIndent)
		for i, p := range n.Params {
			paramStr := p.Name.Lexeme
			if p.Type != nil {
				paramStr += ": " + dump(p.Type, "", true)
			}
			last := i == len(n.Params)-1
			conn := "├─ "
			if last {
				conn = "└─ "
			}
			out += fmt.Sprintf("%s%s%s\n", nextIndent, conn, paramStr)
		}
		if n.Ret != nil {
			out += fmt.Sprintf("%sReturnType:\n%s", nextIndent, dump(n.Ret, nextIndent+"  ", true))
		}
		out += fmt.Sprintf("%sBody:\n%s", nextIndent, dump(n.Body, nextIndent+"  ", true))
		return out
	case *BlockExpr:
		out := fmt.Sprintf("%s%sBlock\n", indent, connector)
		for i, e := range n.Exprs {
			out += dump(e, nextIndent, i == len(n.Exprs)-1) + "\n"
		}
		return strings.TrimRight(out, "\n")
	case *UseExpr:
		out := fmt.Sprintf("%s%sUse\n", indent, connector)
		out += fmt.Sprintf("%sPath: %v\n", nextIndent, n.Path)
		if len(n.Members) > 0 {
			out += fmt.Sprintf("%sMembers: %v\n", nextIndent, n.Members)
		}
		if n.Alias != "" {
			out += fmt.Sprintf("%sAlias: %s\n", nextIndent, n.Alias)
		}
		return out
	case *QualifiedExpr:
		out := fmt.Sprintf("%s%sQualified\n", indent, connector)
		out += fmt.Sprintf("%sLeft:\n%s\n", nextIndent, dump(n.Left, nextIndent+"  ", true))
		out += fmt.Sprintf("%sRight:\n%s", nextIndent, dump(&Identifier{Name: n.Right.Lexeme, Pos: n.Right}, nextIndent+"  ", true))
		return out
	case *FieldAccessExpr:
		out := fmt.Sprintf("%s%sField\n", indent, connector)
		out += fmt.Sprintf("%sLeft:\n%s\n", nextIndent, dump(n.Left, nextIndent+"  ", true))
		out += fmt.Sprintf("%sRight:\n%s", nextIndent, dump(&Identifier{Name: n.Right, Pos: n.Pos}, nextIndent+"  ", true))
		return out
	case *IfExpr:
		out := fmt.Sprintf("%s%sIfExpr\n", indent, connector)
		out += fmt.Sprintf("%sCond:\n%s\n", nextIndent, dump(n.Cond, nextIndent+"  ", true))
		out += fmt.Sprintf("%sThen:\n%s", nextIndent, dump(n.Then, nextIndent+"  ", true))
		if n.Else != nil {
			out += fmt.Sprintf("\n%sElse:\n%s", nextIndent, dump(n.Else, nextIndent+"  ", true))
		}
		return out
	case *MatchExpr:
		out := fmt.Sprintf("%s%sMatchExpr\n", indent, connector)
		out += fmt.Sprintf("%sValue:\n%s\n", nextIndent, dump(n.Value, nextIndent+"  ", true))
		out += fmt.Sprintf("%sArms:\n", nextIndent)
		for _, arm := range n.Arms {
			out += fmt.Sprintf("%s- Pattern:\n%s\n", nextIndent+"  ", dump(arm.Pattern, nextIndent+"    ", true))
			if arm.Guard != nil {
				out += fmt.Sprintf("%s  Guard:\n%s\n", nextIndent+"  ", dump(arm.Guard, nextIndent+"    ", true))
			}
			out += fmt.Sprintf("%s  Body:\n%s\n", nextIndent+"  ", dump(arm.Body, nextIndent+"    ", true))
		}
		return out
	case *PipelineExpr:
		out := fmt.Sprintf("%s%sPipelineExpr\n", indent, connector)
		out += fmt.Sprintf("%sLeft:\n%s\n", nextIndent, dump(n.Left, nextIndent+"  ", true))
		out += fmt.Sprintf("%sRight:\n%s", nextIndent, dump(n.Right, nextIndent+"  ", true))
		return out
	case *ListExpr:
		out := fmt.Sprintf("%s%sList\n", indent, connector)
		nextIndent := indent
		for _, elem := range n.Elements {
			out += fmt.Sprintf("%s%s\n", nextIndent, dump(elem, nextIndent, true))
		}
		return out
	case *TypeExpr:
		if n.Generic != nil {
			return fmt.Sprintf("Type(%s(%s))", n.Name, dump(n.Generic, "", true))
		}
		return fmt.Sprintf("Type(%s)", n.Name)

	case *TupleTypeExpr:
		types := []string{}
		for _, t := range n.Types {
			types = append(types, dump(t, "", true))
		}
		return fmt.Sprintf("TupleType(%s)", strings.Join(types, ", "))
	case *TupleExpr:
		out := fmt.Sprintf("%s%sTuple\n", indent, connector)
		nextIndent := indent + "   "
		for i, elem := range n.Elements {
			last := i == len(n.Elements)-1
			out += dump(elem, nextIndent, last) + "\n"
		}
		return strings.TrimRight(out, "\n")
	case *TypeDeclExpr:
		out := fmt.Sprintf("%s%sTypeDecl(pub=%t, name=%s)\n", indent, connector, n.Pub, n.Name.Lexeme)
		if n.Body != nil {
			out += fmt.Sprintf("%sBody:\n%s", indent+"   ", dump(n.Body, indent+"      ", true))
		}
		return out
	case *RecordTypeExpr:
		out := fmt.Sprintf("%s%sRecordType(%s)\n", indent, connector, n.Name.Lexeme)
		for i, f := range n.Fields {
			last := i == len(n.Fields)-1
			fieldStr := f.Name.Lexeme
			if f.Type != nil {
				fieldStr += ": " + dump(f.Type, "", true)
			}
			out += fmt.Sprintf("%s%s%s\n", indent+"   ", "├─ ", fieldStr)
			if last {
				out = strings.Replace(out, "├─ "+fieldStr, "└─ "+fieldStr, 1)
			}
		}
		return out
	default:
		return fmt.Sprintf("%s%s<unknown %T>", indent, connector, n)
	}
}
