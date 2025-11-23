package parser

import (
	"flint/internal/lexer"
	"fmt"
	"strings"
)

type Node interface {
	NodeType() string
}

type Expr interface {
	Node
	exprNode()
}

type Identifier struct {
	Name string
	Pos  lexer.Token
}

func (i *Identifier) exprNode() {}
func (i *Identifier) NodeType() string {
	return "Identifier"
}

type IntLiteral struct {
	Value int64
	Raw   string
	Pos   lexer.Token
}

func (i *IntLiteral) exprNode() {}
func (i *IntLiteral) NodeType() string {
	return "IntLiteral"
}

type FloatLiteral struct {
	Value float64
	Raw   string
	Pos   lexer.Token
}

func (i *FloatLiteral) exprNode() {}
func (i *FloatLiteral) NodeType() string {
	return "FloatLiteral"
}

type StringLiteral struct {
	Value string
	Pos   lexer.Token
}

func (s *StringLiteral) exprNode() {}
func (s *StringLiteral) NodeType() string {
	return "StringLiteral"
}

type ByteLiteral struct {
	Value byte
	Raw   string
	Pos   lexer.Token
}

func (b *ByteLiteral) exprNode() {}
func (b *ByteLiteral) NodeType() string {
	return "ByteLiteral"
}

type BoolLiteral struct {
	Value bool
	Pos   lexer.Token
}

func (b *BoolLiteral) exprNode() {}
func (b *BoolLiteral) NodeType() string {
	return "BoolLiteral"
}

type PrefixExpr struct {
	Operator lexer.Token
	Right    Expr
}

func (p *PrefixExpr) exprNode() {}
func (p *PrefixExpr) NodeType() string {
	return "PrefixExpr"
}

type InfixExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (i *InfixExpr) exprNode() {}
func (i *InfixExpr) NodeType() string {
	return "InfixExpr"
}

type CallExpr struct {
	Callee Expr
	Args   []Expr
	Pos    lexer.Token
}

func (c *CallExpr) exprNode() {}
func (c *CallExpr) NodeType() string {
	return "CallExpr"
}

type ValDeclExpr struct {
	Name  lexer.Token
	Type  Expr
	Value Expr
}

func (d *ValDeclExpr) exprNode() {}
func (d *ValDeclExpr) NodeType() string {
	return "ValDeclExpr"
}

type MutDeclExpr struct {
	Name  lexer.Token
	Type  Expr
	Value Expr
}

func (d *MutDeclExpr) exprNode() {}
func (d *MutDeclExpr) NodeType() string {
	return "MutDeclExpr"
}

type FuncDeclExpr struct {
	Pub    bool
	Name   lexer.Token
	Params []Param
	Ret    Expr
	Body   Expr
}

func (f *FuncDeclExpr) exprNode() {}
func (f *FuncDeclExpr) NodeType() string {
	return "FuncDeclExpr"
}

type Param struct {
	Name lexer.Token
	Type Expr
}

func (p *Param) exprNode() {}
func (p *Param) NodeType() string {
	return "Param"
}

type BlockExpr struct {
	Exprs []Expr
}

func (b *BlockExpr) exprNode() {}
func (b *BlockExpr) NodeType() string {
	return "BlockExpr"
}

type UseExpr struct {
	Path    []string
	Alias   string
	Members []string
	Pos     lexer.Token
}

func (u *UseExpr) exprNode() {}
func (u *UseExpr) NodeType() string {
	return "UseExpr"
}

type QualifiedExpr struct {
	Left  Expr
	Right lexer.Token
	Pos   lexer.Token
}

func (q *QualifiedExpr) exprNode() {}
func (q *QualifiedExpr) NodeType() string {
	return "UseExpr"
}

type IfExpr struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (i *IfExpr) exprNode() {}
func (i *IfExpr) NodeType() string {
	return "IfExpr"
}

type MatchArm struct {
	Pattern Expr
	Guard   Expr
	Body    Expr
}

func (m *MatchArm) exprNode() {}
func (m *MatchArm) NodeType() string {
	return "MatchArm"
}

type MatchExpr struct {
	Value Expr
	Arms  []*MatchArm
}

func (m *MatchExpr) exprNode() {}
func (m *MatchExpr) NodeType() string {
	return "MatchExpr"
}

type PipelineExpr struct {
	Left  Expr
	Right Expr
}

func (p *PipelineExpr) exprNode() {}
func (p *PipelineExpr) NodeType() string {
	return "PipelineExpr"
}

type ForExpr struct {
	Vars     []Expr
	Iterable Expr
	Where    Expr
	Body     Expr
}

func (f *ForExpr) exprNode() {}
func (f *ForExpr) NodeType() string {
	return "ForExpr"
}

type ListExpr struct {
	Elements []Expr
	Pos      lexer.Token
}

func (l *ListExpr) exprNode() {}
func (l *ListExpr) NodeType() string {
	return "ListExpr"
}

type TypeExpr struct {
	Name    string
	Pos     lexer.Token
	Generic Expr
}

func (t *TypeExpr) exprNode() {}
func (t *TypeExpr) NodeType() string {
	return "TypeExpr"
}

type TupleTypeExpr struct {
	Types []Expr
	Pos   lexer.Token
}

func (t *TupleTypeExpr) exprNode() {}
func (t *TupleTypeExpr) NodeType() string {
	return "TupleTypeExpr"
}

type TupleExpr struct {
	Elements []Expr
	Pos      lexer.Token
}

func (t *TupleExpr) exprNode() {}
func (t *TupleExpr) NodeType() string {
	return "TupleExpr"
}

type RecordTypeExpr struct {
	Name   lexer.Token
	Fields []Param
	Pos    lexer.Token
}

func (r *RecordTypeExpr) exprNode() {}
func (r *RecordTypeExpr) NodeType() string {
	return "RecordTypeExpr"
}

type TypeDeclExpr struct {
	Pub  bool
	Name lexer.Token
	Body Expr
	Pos  lexer.Token
}

func (t *TypeDeclExpr) exprNode() {}
func (t *TypeDeclExpr) NodeType() string {
	return "TypeDeclExpr"
}

type Program struct {
	Exprs []Expr
}

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
		out := fmt.Sprintf("%s%sFuncDecl(pub=%t, name=%s)\n",
			indent, connector, n.Pub, n.Name.Lexeme)
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
	case *ForExpr:
		out := fmt.Sprintf("%s%sForExpr\n", indent, connector)
		nextIndent := indent + "  "

		if len(n.Vars) > 0 {
			out += fmt.Sprintf("%sVars:\n", nextIndent)
			for i, v := range n.Vars {
				conn := "├─ "
				if i == len(n.Vars)-1 {
					conn = "└─ "
				}
				out += fmt.Sprintf("%s%s%s\n", nextIndent+"  ", conn, dump(v, nextIndent+"    ", true))
			}
		}
		if n.Iterable != nil {
			out += fmt.Sprintf("%sIterable:\n%s\n", nextIndent, dump(n.Iterable, nextIndent+"  ", true))
		}
		if n.Where != nil {
			out += fmt.Sprintf("%sWhere:\n%s\n", nextIndent, dump(n.Where, nextIndent+"  ", true))
		}
		if n.Body != nil {
			out += fmt.Sprintf("%sBody:\n%s", nextIndent, dump(n.Body, nextIndent+"  ", true))
		}
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
