package parser

import (
	"flint/internal/lexer"
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

type VarDeclExpr struct {
	Mutable bool
	Name    lexer.Token
	Type    Expr
	Value   Expr
}

func (d *VarDeclExpr) exprNode() {}
func (d *VarDeclExpr) NodeType() string {
	return "VarDeclExpr"
}

type FuncDeclExpr struct {
	Pub        bool
	Recursion  bool
	Name       lexer.Token
	Params     []Param
	Ret        Expr
	Body       Expr
	Decorators []Decorator
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
	return "QualifiedExpr"
}

type FieldAccessExpr struct {
	Left  Expr
	Right string
	Pos   lexer.Token
}

func (q *FieldAccessExpr) exprNode() {}
func (q *FieldAccessExpr) NodeType() string {
	return "FieldAccessExpr"
}

type IfExpr struct {
	Cond Expr
	Then Expr
	Else Expr
	Pos  lexer.Token
}

func (i *IfExpr) exprNode() {}
func (i *IfExpr) NodeType() string {
	return "IfExpr"
}

type MatchArm struct {
	Pattern Expr
	Guard   Expr
	Body    Expr
	Pos     lexer.Token
}

func (m *MatchArm) exprNode() {}
func (m *MatchArm) NodeType() string {
	return "MatchArm"
}

func (m *MatchArm) IsWildCardArm() bool {
	return m.Pattern.NodeType() == "Identifier" && m.Pattern.(*Identifier).Name == "_"
}

type MatchExpr struct {
	Value Expr
	Arms  []*MatchArm
	Pos   lexer.Token
}

func (m *MatchExpr) exprNode() {}
func (m *MatchExpr) NodeType() string {
	return "MatchExpr"
}

type PipelineExpr struct {
	Left  Expr
	Right Expr
	Pos   lexer.Token
}

func (p *PipelineExpr) exprNode() {}
func (p *PipelineExpr) NodeType() string {
	return "PipelineExpr"
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

type Decorator struct {
	Name string
	Args []Expr
	Pos  lexer.Token
}

func (t *Decorator) exprNode() {}
func (t *Decorator) NodeType() string {
	return "Decorator"
}

type AssignExpr struct {
	Name  *Identifier
	Value Expr
	Pos   lexer.Token
}

func (a *AssignExpr) exprNode() {}
func (a *AssignExpr) NodeType() string {
	return "AssignExpr"
}

type Program struct {
	Exprs []Expr
}
