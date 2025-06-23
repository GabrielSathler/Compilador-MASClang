package semantic_analyzer

import (
	"fmt"

	"github.com/GabrielSathler/Compilador-MASClang/ast"
)

type SemanticAnalyzer struct {
	errors []string
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{errors: []string{}}
}

func (s *SemanticAnalyzer) Analyze(node ast.Node) {
	s.analyzeNode(node)
}

func (s *SemanticAnalyzer) analyzeNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Program:
		for _, declaration := range n.Decls {
			s.analyzeNode(declaration)
		}
	case *ast.Function:
		s.analyzeNode(n.Body)
	case *ast.Var:
		if n.Value != nil {
			s.analyzeExpr(n.Value)
		}
	case *ast.CodeBlock:
		for _, stmt := range n.Stmts {
			s.analyzeNode(stmt)
		}
	}
}

func (s *SemanticAnalyzer) analyzeExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.IntLiteral:

	case *ast.Ident:

	case *ast.BinaryExpr:
		s.analyzeExpr(e.Left)
		s.analyzeExpr(e.Right)
	}
}

func (s *SemanticAnalyzer) Errors() []string {
	return s.errors
}

func (s *SemanticAnalyzer) reportError(msg string) {
	s.errors = append(s.errors, msg)
	fmt.Println("Semantic error:", msg)
}
