package semantic_analyzer

import (
	"fmt"

	"github.com/GabrielSathler/Compilador-MASClang/ast"
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type SemanticAnalyzer struct {
	errors []string
	scopes []map[string]string
	funcs  map[string]*ast.Function
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		errors: []string{},
		scopes: []map[string]string{{}}, // escopo global
		funcs:  map[string]*ast.Function{},
	}
}

// ====== APIs externas ======

func (s *SemanticAnalyzer) Analyze(node ast.Node) {
	s.analyzeNode(node)
}

func (s *SemanticAnalyzer) Errors() []string {
	return s.errors
}

// ====== Controle de escopo ======

func (s *SemanticAnalyzer) pushScope() {
	s.scopes = append(s.scopes, map[string]string{})
}

func (s *SemanticAnalyzer) popScope() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *SemanticAnalyzer) declareVar(name, typ string) {
	current := s.scopes[len(s.scopes)-1]
	current[name] = typ
}

func (s *SemanticAnalyzer) lookupVar(name string) (string, bool) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if typ, ok := s.scopes[i][name]; ok {
			return typ, true
		}
	}
	return "", false
}

// ====== Análise de nós ======

func (s *SemanticAnalyzer) analyzeNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Program:
		for _, decl := range n.Decls {
			if fn, ok := decl.(*ast.Function); ok {
				s.funcs[fn.Name] = fn
			}
		}
		for _, decl := range n.Decls {
			s.analyzeNode(decl)
		}

	case *ast.Function:
		s.pushScope()
		for _, param := range n.Params {
			s.declareVar(param.Name, tokens.Token(param.Type).String())
		}
		s.analyzeNode(n.Body)
		s.popScope()

	case *ast.CodeBlock:
		s.pushScope()
		for _, stmt := range n.Stmts {
			s.analyzeNode(stmt)
		}
		s.popScope()

	case *ast.Var:
		typ := tokens.Token(n.Type).String()
		if n.Value != nil {
			valType := s.analyzeExpr(n.Value)
			if valType != typ {
				s.reportError(fmt.Sprintf("Type mismatch in variable '%s': expected %s, got %s", n.Name, typ, valType))
			}
		}
		s.declareVar(n.Name, typ)

	case *ast.Assignment:
		varType, ok := s.lookupVar(n.Name)
		if !ok {
			s.reportError(fmt.Sprintf("Undeclared variable '%s'", n.Name))
			return
		}
		valType := s.analyzeExpr(n.Value)
		if varType != valType {
			s.reportError(fmt.Sprintf("Type mismatch in assignment to '%s': expected %s, got %s", n.Name, varType, valType))
		}

	case *ast.Return:
		if n.Value != nil {
			s.analyzeExpr(n.Value)
		}

	case *ast.If:
		condType := s.analyzeExpr(n.Condition)
		if condType != "bool" {
			s.reportError("Condition in if statement must be boolean")
		}
		s.analyzeNode(n.ThenBlock)
		if n.ElseBlock != nil {
			s.analyzeNode(n.ElseBlock)
		}

	case *ast.While:
		condType := s.analyzeExpr(n.Condition)
		if condType != "bool" {
			s.reportError("Condition in while must be boolean")
		}
		s.analyzeNode(n.Body)

		/*case *ast.For:
			s.pushScope()
			s.analyzeNode(n.Init)
			condType := s.analyzeExpr(n.Condition)
			if condType != "bool" {
				s.reportError("Condition in for must be boolean")
			}
			s.analyzeNode(n.Increment)
			s.analyzeNode(n.Body)
			s.popScope()

		case *ast.Print:
			s.analyzeExpr(n.Value)

		case *ast.Input:
			_, ok := s.lookupVar(n.Name)
			if !ok {
				s.reportError(fmt.Sprintf("Undeclared variable '%s' in input", n.Name))
			}*/
	}
}

// ====== Análise de expressões ======

func (s *SemanticAnalyzer) analyzeExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.IntLiteral:
		return "int"

	case *ast.FloatLiteral:
		return "float"

	case *ast.StringLiteral:
		return "string"

	case *ast.CharLiteral:
		return "char"

	case *ast.BoolLiteral:
		return "bool"

	case *ast.Ident:
		typ, ok := s.lookupVar(e.Name)
		if !ok {
			s.reportError(fmt.Sprintf("Undeclared variable '%s'", e.Name))
			return "unknown"
		}
		return typ

	case *ast.BinaryExpr:
		leftType := s.analyzeExpr(e.Left)
		rightType := s.analyzeExpr(e.Right)

		if isArithmeticOp(e.Op) {
			if leftType != "int" && leftType != "float" {
				s.reportError(fmt.Sprintf("Invalid left operand type %s for arithmetic operator", leftType))
			}
			if rightType != leftType {
				s.reportError(fmt.Sprintf("Type mismatch in binary expression: %s vs %s", leftType, rightType))
			}
			return leftType
		}

		if isComparisonOp(e.Op) {
			if leftType != rightType {
				s.reportError(fmt.Sprintf("Type mismatch in comparison: %s vs %s", leftType, rightType))
			}
			return "bool"
		}

		if e.Op == tokens.DOT {
			if leftType != "string" && rightType != "string" {
				s.reportError("Both operands of '.' must be string")
			}
			return "string"
		}

		s.reportError("Unknown binary operator")
		return "unknown"

	case *ast.FuncCall:
		fn, ok := s.funcs[e.Name]
		if !ok {
			s.reportError(fmt.Sprintf("Undefined function '%s'", e.Name))
			return "unknown"
		}

		if len(fn.Params) != len(e.Args) {
			s.reportError(fmt.Sprintf("Argument count mismatch in function '%s'", e.Name))
		} else {
			for i, param := range fn.Params {
				argType := s.analyzeExpr(e.Args[i])
				paramType := tokens.Token(param.Type).String()
				if argType != paramType {
					s.reportError(fmt.Sprintf("Type mismatch in argument %d of function '%s': expected %s, got %s",
						i+1, e.Name, paramType, argType))
				}
			}
		}
		return tokens.Token(fn.ReturnType).String()

	default:
		s.reportError("Unknown expression type")
		return "unknown"
	}
}

// ====== Operadores ======

func isArithmeticOp(op tokens.Token) bool {
	return op == tokens.ADD || op == tokens.SUB || op == tokens.MUL ||
		op == tokens.DIV || op == tokens.REM
}

func isComparisonOp(op tokens.Token) bool {
	return op == tokens.EQUAL || op == tokens.NEQUAL || op == tokens.LT ||
		op == tokens.LTOE || op == tokens.GT || op == tokens.GTOE
}

// ====== Erro ======

func (s *SemanticAnalyzer) reportError(msg string) {
	s.errors = append(s.errors, msg)
}
