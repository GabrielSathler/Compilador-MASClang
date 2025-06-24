package syntactic_analyzer

import (
	"fmt"
	"io"
	"strconv"

	"github.com/GabrielSathler/Compilador-MASClang/ast"
	"github.com/GabrielSathler/Compilador-MASClang/lexical_analyzer"
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type Parser struct {
	lexer     *lexical_analyzer.Lexer
	currToken tokens.Token
	currLex   string
	pos       lexical_analyzer.Position
}

func NewParser(reader io.Reader) *Parser {
	lexer := lexical_analyzer.NewLexer(reader)
	p := &Parser{lexer: lexer}
	p.advance()

	return p
}

func (p *Parser) advance() {
	p.pos, p.currToken, p.currLex = p.lexer.Lex()
}

func (p *Parser) expect(expectedToken tokens.Token) {
	if p.currToken != expectedToken {
		panic(fmt.Sprintf("expected token %v, got token: %v at position: %v", expectedToken, p.currToken, p.pos))
	}

	p.advance()
}

func isValidType(t tokens.Token) bool {
	return t == tokens.INT || t == tokens.STRING || t == tokens.FLOAT || t == tokens.CHAR || t == tokens.BOOL
}

func (p *Parser) parseBlock() *ast.CodeBlock {
	p.expect(tokens.LBRACE)
	statements := []ast.Node{}

	for p.currToken != tokens.RBRACE && p.currToken != tokens.EOF {
		statement := p.parseStatement()

		if statement != nil {
			statements = append(statements, statement)
		}
	}

	p.expect(tokens.RBRACE)

	return &ast.CodeBlock{Statements: statements}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Declarations: []ast.Node{}}

	for p.currToken != tokens.EOF {
		switch p.currToken {
		case tokens.FUNC:
			fn := p.parseFunction()
			program.Declarations = append(program.Declarations, fn)
		case tokens.VAR:
			v := p.parseVar()
			program.Declarations = append(program.Declarations, v)
		case tokens.IF:
			statement := p.parseIf()
			program.Declarations = append(program.Declarations, statement)
		case tokens.FOR:
			statement := p.parseFor()
			program.Declarations = append(program.Declarations, statement)
		case tokens.WHILE:
			statement := p.parseWhile()
			program.Declarations = append(program.Declarations, statement)
		case tokens.PRINT:
			statement := p.parsePrint()
			program.Declarations = append(program.Declarations, statement)
		case tokens.INPUT:
			statement := p.parseInput()
			program.Declarations = append(program.Declarations, statement)
		case tokens.RETURN:
			statement := p.parseReturn()
			program.Declarations = append(program.Declarations, statement)
		case tokens.IDENT:
			statement := p.parseAssignmentOrFuncCall(true)
			program.Declarations = append(program.Declarations, statement)
		default:
			panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
		}
	}

	return program
}

func (p *Parser) parseFunction() *ast.Function {
	p.expect(tokens.FUNC)

	name := p.currLex

	p.expect(tokens.IDENT)
	p.expect(tokens.LPAREN)

	params := p.parseFunctionParameters()

	p.expect(tokens.RPAREN)
	p.expect(tokens.COLON)

	returnType := p.currToken

	p.advance()
	body := p.parseBlock()

	return &ast.Function{Name: name, Params: params, ReturnType: returnType, Body: body}
}

func (p *Parser) parseFunctionParameters() []ast.Param {
	params := []ast.Param{}

	for p.currToken != tokens.RPAREN {
		name := p.currLex

		p.expect(tokens.IDENT)
		p.expect(tokens.COLON)

		parameterType := p.currToken

		p.advance()
		params = append(params, ast.Param{Name: name, Type: parameterType})

		if p.currToken == tokens.COMMA {
			p.advance()
		}
	}

	return params
}

func (p *Parser) parseStatement() ast.Node {
	switch p.currToken {
	case tokens.IF:
		return p.parseIf()
	case tokens.VAR:
		return p.parseVar()
	case tokens.FOR:
		return p.parseFor()
	case tokens.WHILE:
		return p.parseWhile()
	case tokens.PRINT:
		return p.parsePrint()
	case tokens.INPUT:
		return p.parseInput()
	case tokens.RETURN:
		return p.parseReturn()
	case tokens.IDENT:
		return p.parseAssignmentOrFuncCall(true)
	default:
		panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
	}
}

func (p *Parser) parseIf() ast.Node {
	p.expect(tokens.IF)
	p.expect(tokens.LPAREN)

	condition := p.parseComparison()

	p.expect(tokens.RPAREN)

	thenBlock := p.parseBlock()

	var elseBlock *ast.CodeBlock = nil
	if p.currToken == tokens.ELSE {
		p.advance()
		elseBlock = p.parseBlock()
	}

	return &ast.If{Condition: condition, ThenBlock: thenBlock, ElseBlock: elseBlock}
}

func (p *Parser) parseVar() *ast.Var {
	p.expect(tokens.VAR)

	if p.currToken != tokens.IDENT {
		panic(fmt.Sprintf("expected variable name, got %v at %v", p.currToken, p.pos))
	}

	name := p.currLex
	p.advance()

	p.expect(tokens.COLON)

	if !isValidType(p.currToken) {
		panic(fmt.Sprintf("expected variable type, got %v at %v", p.currToken, p.pos))
	}

	typeTok := p.currToken
	p.advance()

	var value ast.Expression = nil
	if p.currToken == tokens.ASSIGN {
		p.advance()
		value = p.parseComparison()
	}

	p.expect(tokens.SEMI)

	return &ast.Var{Name: name, Type: typeTok, Value: value}
}

func (p *Parser) parseFor() ast.Node {
	p.expect(tokens.FOR)
	p.expect(tokens.LPAREN)

	var init ast.Node = nil
	if p.currToken == tokens.VAR {
		init = p.parseVar()
	} else {
		init = p.parseAssignmentOrFuncCall(true)
	}

	condition := p.parseComparison()

	p.expect(tokens.SEMI)

	post := p.parseAssignmentOrFuncCall(false)

	p.expect(tokens.RPAREN)

	body := p.parseBlock()
	return &ast.For{Init: init, Condition: condition, Increment: post, Body: body}
}

func (p *Parser) parseWhile() ast.Node {
	p.expect(tokens.WHILE)
	p.expect(tokens.LPAREN)

	condition := p.parseComparison()

	p.expect(tokens.RPAREN)

	body := p.parseBlock()
	return &ast.While{Condition: condition, Body: body}
}

func (p *Parser) parsePrint() ast.Node {
	p.expect(tokens.PRINT)
	p.expect(tokens.LPAREN)

	value := p.parseAdditive()

	p.expect(tokens.RPAREN)
	p.expect(tokens.SEMI)

	return &ast.Print{Value: value}
}

func (p *Parser) parseInput() ast.Node {
	p.expect(tokens.INPUT)
	p.expect(tokens.LPAREN)

	value := p.currLex

	p.expect(tokens.IDENT)
	p.expect(tokens.RPAREN)
	p.expect(tokens.SEMI)

	return &ast.Input{Value: value}
}

func (p *Parser) parseReturn() ast.Node {
	p.expect(tokens.RETURN)

	var value ast.Expression = nil
	if p.currToken != tokens.SEMI {
		value = p.parseAdditive()
	}

	p.expect(tokens.SEMI)

	return &ast.Return{Value: value}
}

func (p *Parser) parseAssignmentOrFuncCall(requireSemi bool) ast.Node {
	if p.currToken != tokens.IDENT {
		panic(fmt.Sprintf("expected identifier, got %v at %v", p.currToken, p.pos))
	}

	name := p.currLex
	p.advance()

	switch p.currToken {
	case tokens.ASSIGN:
		p.advance()
		value := p.parseAdditive()

		if requireSemi {
			if p.currToken != tokens.SEMI {
				panic(fmt.Sprintf("expected token ;, got token: %v at position: %v", p.currToken, p.pos))
			}

			p.advance()
		}

		return &ast.Assign{Name: name, Value: value}
	case tokens.LPAREN:
		p.advance()
		arguments := []ast.Expression{}

		if p.currToken != tokens.RPAREN {
			for {
				argument := p.parseAdditive()
				arguments = append(arguments, argument)

				if p.currToken == tokens.COMMA {
					p.advance()
					continue
				}

				break
			}
		}

		p.expect(tokens.RPAREN)

		if requireSemi {
			p.expect(tokens.SEMI)
		}

		return &ast.FuncCall{Name: name, Arguments: arguments}
	default:
		panic(fmt.Sprintf("unexpected token after identifier %v at %v", p.currToken, p.pos))
	}
}

func (p *Parser) parseComparison() ast.Expression {
	left := p.parseAdditive()

	for p.currToken == tokens.EQUAL || p.currToken == tokens.NEQUAL ||
		p.currToken == tokens.LT || p.currToken == tokens.LTOE ||
		p.currToken == tokens.GT || p.currToken == tokens.GTOE {
		operation := p.currToken
		p.advance()

		right := p.parseAdditive()
		left = &ast.BinaryExpression{Left: left, Operation: operation, Right: right}
	}

	return left
}

func (p *Parser) parseAdditive() ast.Expression {
	left := p.parseMultiplicative()

	for p.currToken == tokens.ADD || p.currToken == tokens.SUB || p.currToken == tokens.DOT {
		operation := p.currToken
		p.advance()

		right := p.parseMultiplicative()
		left = &ast.BinaryExpression{Left: left, Operation: operation, Right: right}
	}

	return left
}

func (p *Parser) parseMultiplicative() ast.Expression {
	left := p.parseFactor()

	for p.currToken == tokens.MUL || p.currToken == tokens.DIV || p.currToken == tokens.REM {
		operation := p.currToken
		p.advance()

		right := p.parseFactor()
		left = &ast.BinaryExpression{Left: left, Operation: operation, Right: right}
	}

	return left
}

func (p *Parser) parseFactor() ast.Expression {
	switch p.currToken {
	case tokens.INT:
		stringValue := p.currLex
		p.advance()

		value, err := strconv.Atoi(stringValue)

		if err != nil {
			panic(fmt.Sprintf("invalid integer literal: %v", stringValue))
		}

		return &ast.IntLiteral{Value: value}
	case tokens.STRING:
		value := p.currLex
		p.advance()

		return &ast.StringLiteral{Value: value}
	case tokens.CHAR:
		val := p.currLex
		p.advance()

		if len(val) == 3 && val[0] == '\'' && val[2] == '\'' {
			return &ast.CharLiteral{Value: rune(val[1])}
		} else if len(val) == 1 {
			return &ast.CharLiteral{Value: rune(val[0])}
		} else {
			panic(fmt.Sprintf("invalid char literal: %v", val))
		}
	case tokens.FLOAT:
		stringValue := p.currLex
		p.advance()

		value, err := strconv.ParseFloat(stringValue, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid float literal: %v", stringValue))
		}

		return &ast.FloatLiteral{Value: value}
	case tokens.TRUE, tokens.FALSE:
		value := (p.currToken == tokens.TRUE)
		p.advance()

		return &ast.BoolLiteral{Value: value}
	case tokens.IDENT:
		name := p.currLex
		p.advance()

		if p.currToken == tokens.LPAREN {
			p.advance()
			arguments := []ast.Expression{}

			if p.currToken != tokens.RPAREN {
				for {
					arguments = append(arguments, p.parseComparison())
					if p.currToken == tokens.COMMA {
						p.advance()
					} else {
						break
					}
				}
			}

			p.expect(tokens.RPAREN)
			return &ast.FuncCall{Name: name, Arguments: arguments}
		}

		return &ast.Ident{Name: name}
	default:
		panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
	}
}
