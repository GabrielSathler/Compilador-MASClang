package ast

import (
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type Node interface {
	Pos() int
	Line() int
}

type Expression interface {
	Node
}

type Program struct {
	Declarations []Node
	LineIdent    int
}

func (p *Program) Pos() int  { return 0 }
func (p *Program) Line() int { return p.LineIdent }

type Function struct {
	Name       string
	Params     []Param
	ReturnType tokens.Token
	Body       *CodeBlock
	LineIdent  int
}

func (f *Function) Pos() int  { return 0 }
func (f *Function) Line() int { return f.LineIdent }

type Param struct {
	Name string
	Type tokens.Token
}

type CodeBlock struct {
	Statements []Node
	LineIdent  int
}

func (b *CodeBlock) Pos() int  { return 0 }
func (b *CodeBlock) Line() int { return b.LineIdent }

type Var struct {
	Name      string
	Type      tokens.Token
	Value     Expression
	LineIdent int
}

func (v *Var) Pos() int  { return 0 }
func (v *Var) Line() int { return v.LineIdent }

type Assignment struct {
	Name      string
	Value     Expression
	LineIdent int
}

func (a *Assignment) Pos() int  { return 0 }
func (a *Assignment) Line() int { return a.LineIdent }

type Return struct {
	Value     Expression
	LineIdent int
}

func (r *Return) Pos() int  { return 0 }
func (r *Return) Line() int { return r.LineIdent }

type IntLiteral struct {
	Value     int
	LineIdent int
}

func (i *IntLiteral) Pos() int  { return 0 }
func (i *IntLiteral) Line() int { return i.LineIdent }

type FloatLiteral struct {
	Value     float64
	LineIdent int
}

func (f *FloatLiteral) Pos() int  { return 0 }
func (f *FloatLiteral) Line() int { return f.LineIdent }

type StringLiteral struct {
	Value     string
	LineIdent int
}

func (s *StringLiteral) Pos() int  { return 0 }
func (s *StringLiteral) Line() int { return s.LineIdent }

type CharLiteral struct {
	Value     rune
	LineIdent int
}

func (c *CharLiteral) Pos() int  { return 0 }
func (c *CharLiteral) Line() int { return c.LineIdent }

type BoolLiteral struct {
	Value     bool
	LineIdent int
}

func (b *BoolLiteral) Pos() int  { return 0 }
func (b *BoolLiteral) Line() int { return b.LineIdent }

type Ident struct {
	Name      string
	LineIdent int
}

func (i *Ident) Pos() int  { return 0 }
func (i *Ident) Line() int { return i.LineIdent }

type BinaryExpression struct {
	Left      Expression
	Operation tokens.Token
	Right     Expression
	LineIdent int
}

func (b *BinaryExpression) Pos() int  { return 0 }
func (b *BinaryExpression) Line() int { return b.LineIdent }

type If struct {
	Condition Expression
	ThenBlock *CodeBlock
	ElseBlock *CodeBlock
	LineIdent int
}

func (i *If) Pos() int  { return 0 }
func (i *If) Line() int { return i.LineIdent }

type For struct {
	Init      Node
	Condition Expression
	Increment Node
	Body      *CodeBlock
	LineIdent int
}

func (f *For) Pos() int  { return 0 }
func (f *For) Line() int { return f.LineIdent }

type While struct {
	Condition Expression
	Body      *CodeBlock
	LineIdent int
}

func (w *While) Pos() int  { return 0 }
func (w *While) Line() int { return w.LineIdent }

type Print struct {
	Value     Expression
	LineIdent int
}

func (p *Print) Pos() int  { return 0 }
func (p *Print) Line() int { return p.LineIdent }

type Input struct {
	Value     string
	LineIdent int
}

func (i *Input) Pos() int  { return 0 }
func (i *Input) Line() int { return i.LineIdent }

type Assign struct {
	Name      string
	Value     Expression
	LineIdent int
}

func (a *Assign) Pos() int  { return 0 }
func (a *Assign) Line() int { return a.LineIdent }

type FuncCall struct {
	Name      string
	Arguments []Expression
	LineIdent int
}

func (f *FuncCall) Pos() int  { return 0 }
func (f *FuncCall) Line() int { return f.LineIdent }
