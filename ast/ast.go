package ast

import (
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type Node interface {
	Pos() int
}

type Expression interface {
	Node
}

type Program struct {
	Declarations []Node
}

func (p *Program) Pos() int { return 0 }

type Function struct {
	Name       string
	Params     []Param
	ReturnType tokens.Token
	Body       *CodeBlock
}

func (f *Function) Pos() int { return 0 }

type Param struct {
	Name string
	Type tokens.Token
}

type CodeBlock struct {
	Statements []Node
}

func (b *CodeBlock) Pos() int { return 0 }

type Var struct {
	Name  string
	Type  tokens.Token
	Value Expression
}

func (v *Var) Pos() int { return 0 }

type Assignment struct {
	Name  string
	Value Expression
}

func (a *Assignment) Pos() int { return 0 }

type Return struct {
	Value Expression
}

func (r *Return) Pos() int { return 0 }

type IntLiteral struct {
	Value int
}

func (i *IntLiteral) Pos() int { return 0 }

type FloatLiteral struct {
	Value float64
}

func (f *FloatLiteral) Pos() int { return 0 }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Pos() int { return 0 }

type CharLiteral struct {
	Value rune
}

func (c *CharLiteral) Pos() int { return 0 }

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) Pos() int { return 0 }

type Ident struct {
	Name string
}

func (i *Ident) Pos() int { return 0 }

type BinaryExpression struct {
	Left      Expression
	Operation tokens.Token
	Right     Expression
}

func (b *BinaryExpression) Pos() int { return 0 }

type If struct {
	Condition Expression
	ThenBlock *CodeBlock
	ElseBlock *CodeBlock
}

func (i *If) Pos() int { return 0 }

type For struct {
	Init      Node
	Condition Expression
	Increment Node
	Body      *CodeBlock
}

func (f *For) Pos() int { return 0 }

type While struct {
	Condition Expression
	Body      *CodeBlock
}

func (w *While) Pos() int { return 0 }

type Print struct {
	Value Expression
}

func (p *Print) Pos() int { return 0 }

type Input struct {
	Value string
}

func (i *Input) Pos() int { return 0 }

type Assign struct {
	Name  string
	Value Expression
}

func (a *Assign) Pos() int { return 0 }

type FuncCall struct {
	Name      string
	Arguments []Expression
}

func (f *FuncCall) Pos() int { return 0 }
