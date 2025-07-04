package semantic_analyzer

import (
	"fmt"
	"strconv"

	"github.com/GabrielSathler/Compilador-MASClang/ast"
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type SemanticAnalyzer struct {
	Errors []string
	scopes []map[string]string
	funcs  map[string]*ast.Function
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		Errors: []string{},
		scopes: []map[string]string{{}},
		funcs:  map[string]*ast.Function{},
	}
}

func (s *SemanticAnalyzer) Analyze(node ast.Node) {
	s.analyzeNode(node)
}

func (s *SemanticAnalyzer) analyzeNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Program:
		for _, declaration := range n.Declarations {
			if function, ok := declaration.(*ast.Function); ok {
				s.funcs[function.Name] = function
			}
		}

		for _, declaration := range n.Declarations {
			s.analyzeNode(declaration)
		}
	case *ast.Function:
		s.scopes = append(s.scopes, map[string]string{})

		for _, param := range n.Params {
			s.declareVar(param.Name, tokens.Token(param.Type).String())
		}

		s.analyzeNode(n.Body)

		s.scopes = s.scopes[:len(s.scopes)-1]
	case *ast.CodeBlock:
		s.scopes = append(s.scopes, map[string]string{})

		for _, stmt := range n.Statements {
			s.analyzeNode(stmt)
		}

		s.scopes = s.scopes[:len(s.scopes)-1]
	case *ast.Var:
		varType := tokens.Token(n.Type).String()

		if n.Value != nil {
			valueType := s.analyzeExpression(n.Value)

			if valueType != varType {
				s.reportError(fmt.Sprintf("type mismatch in variable '%s': expected %s, got %s at line %s", n.Name, varType, valueType, strconv.Itoa(n.LineIdent)))
			}
		}

		s.declareVar(n.Name, varType)
	case *ast.Assignment:
		varType, ok := s.lookupVar(n.Name)
		if !ok {
			s.reportError(fmt.Sprintf("undeclared variable '%s' at line %s", n.Name, strconv.Itoa(n.LineIdent)))
			return
		}

		valueType := s.analyzeExpression(n.Value)
		if varType != valueType {
			s.reportError(fmt.Sprintf("type mismatch in assignment to '%s': expected %s, got %s at line %s", n.Name, varType, valueType, strconv.Itoa(n.LineIdent)))
		}
	case *ast.Return:
		if n.Value != nil {
			s.analyzeExpression(n.Value)
		}
	case *ast.If:
		condition := s.analyzeExpression(n.Condition)
		if condition != "bool" {
			s.reportError(fmt.Sprintf("condition in if statement must be boolean at line %s", strconv.Itoa(n.LineIdent)))
		}

		s.analyzeNode(n.ThenBlock)
		if n.ElseBlock != nil {
			s.analyzeNode(n.ElseBlock)
		}
	case *ast.While:
		condition := s.analyzeExpression(n.Condition)
		if condition != "bool" {
			s.reportError(fmt.Sprintf("condition in while must be boolean at line %s", strconv.Itoa(n.LineIdent)))
		}

		s.analyzeNode(n.Body)

	case *ast.For:
		s.scopes = append(s.scopes, map[string]string{})
		s.analyzeNode(n.Init)

		condition := s.analyzeExpression(n.Condition)
		if condition != "bool" {
			s.reportError(fmt.Sprintf("condition in for must be boolean at line %s", strconv.Itoa(n.LineIdent)))
		}

		s.analyzeNode(n.Increment)
		s.analyzeNode(n.Body)
		s.scopes = s.scopes[:len(s.scopes)-1]
	case *ast.Print:
		s.analyzeExpression(n.Value)
	case *ast.Input:
		_, ok := s.lookupVar(n.Value)
		if !ok {
			s.reportError(fmt.Sprintf("undeclared variable '%s' in input at line %s", n.Value, strconv.Itoa(n.LineIdent)))
		}
	}
}

func (s *SemanticAnalyzer) declareVar(name, varType string) {
	current := s.scopes[len(s.scopes)-1]
	current[name] = varType
}

func (s *SemanticAnalyzer) lookupVar(name string) (string, bool) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if varType, ok := s.scopes[i][name]; ok {
			return varType, true
		}
	}

	return "", false
}

func (s *SemanticAnalyzer) analyzeExpression(expression ast.Expression) string {
	switch e := expression.(type) {
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
		varType, ok := s.lookupVar(e.Name)

		if !ok {
			s.reportError(fmt.Sprintf("undeclared variable '%s' at line %s", e.Name, strconv.Itoa(e.LineIdent)))
			return "unknown"
		}

		return varType
	case *ast.BinaryExpression:
		leftType := s.analyzeExpression(e.Left)
		rightType := s.analyzeExpression(e.Right)

		if e.Operation == tokens.ADD || e.Operation == tokens.DOT {
			if leftType == "string" || rightType == "string" {
				return "string"
			}

			if leftType == "int" && rightType == "int" {
				return "int"
			}

			if leftType == "float" && rightType == "float" {
				return "float"
			}

			s.reportError(fmt.Sprintf("invalid operand types for '+' at line %s", strconv.Itoa(e.LineIdent)))
			return "unknown"
		}

		if isArithmeticOperation(e.Operation) {
			if leftType != "int" && leftType != "float" {
				s.reportError(fmt.Sprintf("invalid left operand type %s for arithmetic operator at line %s", leftType, strconv.Itoa(e.LineIdent)))
			}

			if rightType != leftType {
				s.reportError(fmt.Sprintf("type mismatch in binary expression: %s vs %s at line %s", leftType, rightType, strconv.Itoa(e.LineIdent)))
			}

			return leftType
		}

		if isComparisonOperation(e.Operation) {
			if leftType != rightType {
				s.reportError(fmt.Sprintf("type mismatch in comparison: %s vs %s at line %s", leftType, rightType, strconv.Itoa(e.LineIdent)))
			}

			return "bool"
		}

		s.reportError(fmt.Sprintf("unknown binary operator at line %s", strconv.Itoa(e.LineIdent)))
		return "unknown"
	case *ast.FuncCall:
		fn, ok := s.funcs[e.Name]
		if !ok {
			s.reportError(fmt.Sprintf("undefined function '%s' at line %s", e.Name, strconv.Itoa(e.LineIdent)))
			return "unknown"
		}

		if len(fn.Params) != len(e.Arguments) {
			s.reportError(fmt.Sprintf("argument count mismatch in function '%s' at line %s", e.Name, strconv.Itoa(e.LineIdent)))
		} else {
			for i, param := range fn.Params {
				argumentType := s.analyzeExpression(e.Arguments[i])
				paramType := tokens.Token(param.Type).String()

				if argumentType != paramType {
					s.reportError(fmt.Sprintf(
						"type mismatch in argument %d of function '%s': expected %s, got %s at line %s",
						i+1,
						e.Name,
						paramType,
						argumentType,
						strconv.Itoa(e.LineIdent),
					))
				}
			}
		}

		return tokens.Token(fn.ReturnType).String()
	default:
		s.reportError(fmt.Sprintf("unknown expression type at line %s", strconv.Itoa(expression.Line())))
		return "unknown"
	}
}

func isArithmeticOperation(operation tokens.Token) bool {
	return operation == tokens.ADD || operation == tokens.SUB || operation == tokens.MUL ||
		operation == tokens.DIV || operation == tokens.REM
}

func isComparisonOperation(operation tokens.Token) bool {
	return operation == tokens.EQUAL || operation == tokens.NEQUAL || operation == tokens.LT ||
		operation == tokens.LTOE || operation == tokens.GT || operation == tokens.GTOE
}

func (s *SemanticAnalyzer) reportError(msg string) {
	s.Errors = append(s.Errors, msg)
}
