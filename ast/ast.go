package ast

import (
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

// ===== Interfaces =====

type Node interface {
	Pos() int
}

type Expr interface {
	Node
}

// ===== Programa =====

type Program struct {
	Decls []Node
}

func (p *Program) Pos() int { return 0 }

// ===== Função =====

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

// ===== Bloco de Código =====

type CodeBlock struct {
	Stmts []Node
}

func (b *CodeBlock) Pos() int { return 0 }

// ===== Declaração de Variável =====

type Var struct {
	Name  string
	Type  tokens.Token
	Value Expr
}

func (v *Var) Pos() int { return 0 }

// ===== Atribuição =====

type Assignment struct {
	Name  string
	Value Expr
}

func (a *Assignment) Pos() int { return 0 }

// ===== Retorno =====

type Return struct {
	Value Expr // pode ser nil para retorno vazio
}

func (r *Return) Pos() int { return 0 }

// ===== Expressões Literais =====

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

// ===== Identificador =====

type Ident struct {
	Name string
}

func (i *Ident) Pos() int { return 0 }

// ===== Operações Binárias =====

type BinaryExpr struct {
	Left  Expr
	Op    tokens.Token
	Right Expr
}

func (b *BinaryExpr) Pos() int { return 0 }

// ===== Chamada de Função =====

type FuncCall struct {
	Name string
	Args []Expr
}

func (f *FuncCall) Pos() int { return 0 }

// ===== Condicional If =====

type If struct {
	Condition Expr
	ThenBlock *CodeBlock
	ElseBlock *CodeBlock // Pode ser nil
}

func (i *If) Pos() int { return 0 }

// ===== Laço While =====

type While struct {
	Condition Expr
	Body      *CodeBlock
}

func (w *While) Pos() int { return 0 }

// ===== Laço For =====

type For struct {
	Init      Node
	Condition Expr
	Increment Node
}
