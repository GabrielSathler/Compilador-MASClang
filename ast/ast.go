package ast

import (
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type Node interface {
	Pos() int
}

type Program struct {
	Decls []Node
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
	Stmts []Node
}

func (b *CodeBlock) Pos() int { return 0 }

type Var struct {
	Name  string
	Type  tokens.Token
	Value Expr
}

func (v *Var) Pos() int { return 0 }

type Expr interface {
	Node
}

type IntLiteral struct {
	Value int
}

func (i *IntLiteral) Pos() int { return 0 }

type Ident struct {
	Name string
}

func (i *Ident) Pos() int { return 0 }

type BinaryExpr struct {
	Left  Expr
	Op    tokens.Token
	Right Expr
}

func (b *BinaryExpr) Pos() int { return 0 }
