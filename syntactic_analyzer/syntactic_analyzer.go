package syntactic_analyzer

import (
	"fmt"
	"io"

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

func (p *Parser) parseBlock() {
	p.expect(tokens.LBRACE)

	for p.currToken != tokens.RBRACE && p.currToken != tokens.EOF {
		p.parseStatement()
	}

	p.expect(tokens.RBRACE)
}

func isValidType(t tokens.Token) bool {
	return t == tokens.INT || t == tokens.STRING || t == tokens.FLOAT || t == tokens.CHAR || t == tokens.BOOL
}

func (p *Parser) ParseProgram() {
	for p.currToken != tokens.EOF {
		switch p.currToken {
		case tokens.FUNC:
			p.parseFunction()
		case tokens.IF:
			p.parseIf()
		case tokens.VAR:
			p.parseVar()
		case tokens.FOR:
			p.parseFor()
		case tokens.WHILE:
			p.parseWhile()
		case tokens.PRINT:
			p.parsePrint()
		case tokens.INPUT:
			p.parseInput()
		case tokens.IDENT:
			p.parseAssignmentOrFuncCall(true)
		default:
			panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
		}
	}
}

func (p *Parser) parseStatement() {
	switch p.currToken {
	case tokens.IF:
		p.parseIf()
	case tokens.VAR:
		p.parseVar()
	case tokens.FOR:
		p.parseFor()
	case tokens.WHILE:
		p.parseWhile()
	case tokens.PRINT:
		p.parsePrint()
	case tokens.INPUT:
		p.parseInput()
	case tokens.RETURN:
		p.parseReturn()
	case tokens.IDENT:
		p.parseAssignmentOrFuncCall(true)
	default:
		panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
	}
}

func (p *Parser) parseAssignmentOrFuncCall(requireSemi bool) {
	if p.currToken != tokens.IDENT {
		panic(fmt.Sprintf("expected identifier, got %v at %v", p.currToken, p.pos))
	}

	p.advance()

	switch p.currToken {
	case tokens.ASSIGN:
		p.advance()
		p.parseAdditive()

		if requireSemi {
			if p.currToken != tokens.SEMI {
				panic(fmt.Sprintf("expected token ;, got token: %v at position: %v", p.currToken, p.pos))
			}

			p.advance()
		}
	case tokens.LPAREN:
		p.advance()
		p.parseArguments()
		p.expect(tokens.RPAREN)

		if requireSemi {
			p.expect(tokens.SEMI)
		}
	default:
		panic(fmt.Sprintf("unexpected token after identifier %v at %v", p.currToken, p.pos))
	}
}

func (p *Parser) parseFunction() {
	p.expect(tokens.FUNC)

	if p.currToken != tokens.IDENT {
		panic(fmt.Sprintf("expected function name, got %v at %v", p.currToken, p.pos))
	}

	p.advance()
	p.expect(tokens.LPAREN)
	p.parseFunctionParameters()
	p.expect(tokens.RPAREN)
	p.expect(tokens.COLON)

	if !isValidType(p.currToken) {
		panic(fmt.Sprintf("expected return type, got %v at %v", p.currToken, p.pos))
	}

	p.advance()
	p.parseBlock()
}

func (p *Parser) parseFunctionParameters() {
	for p.currToken != tokens.RPAREN {
		if p.currToken != tokens.IDENT {
			panic(fmt.Sprintf("expected parameter name, got %v at %v", p.currToken, p.pos))
		}

		p.advance()
		p.expect(tokens.COLON)

		if !isValidType(p.currToken) {
			panic(fmt.Sprintf("expected parameter type, got %v at %v", p.currToken, p.pos))
		}

		p.advance()

		if p.currToken == tokens.COMMA {
			p.advance()
		}
	}
}

func (p *Parser) parseIf() {
	p.expect(tokens.IF)
	p.expect(tokens.LPAREN)
	p.parseComparison()
	p.expect(tokens.RPAREN)
	p.parseBlock()

	if p.currToken == tokens.ELSE {
		p.advance()
		p.parseBlock()
	}
}

func (p *Parser) parseWhile() {
	p.expect(tokens.WHILE)
	p.expect(tokens.LPAREN)
	p.parseComparison()
	p.expect(tokens.RPAREN)
	p.parseBlock()
}

func (p *Parser) parseVar() {
	p.expect(tokens.VAR)

	if p.currToken != tokens.IDENT {
		panic(fmt.Sprintf("expected variable name, got %v at %v", p.currToken, p.pos))
	}

	p.advance()
	p.expect(tokens.COLON)

	if !isValidType(p.currToken) {
		panic(fmt.Sprintf("expected variable type, got %v at %v", p.currToken, p.pos))
	}

	p.advance()

	if p.currToken == tokens.ASSIGN {
		p.advance()
		p.parseAdditive()
	}

	p.expect(tokens.SEMI)
}

func (p *Parser) parsePrint() {
	p.expect(tokens.PRINT)
	p.expect(tokens.LPAREN)
	p.parseAdditive()
	p.expect(tokens.RPAREN)
	p.expect(tokens.SEMI)
}

func (p *Parser) parseInput() {
	p.expect(tokens.INPUT)
	p.expect(tokens.LPAREN)
	p.expect(tokens.IDENT)
	p.expect(tokens.RPAREN)
	p.expect(tokens.SEMI)
}

func (p *Parser) parseReturn() {
	p.expect(tokens.RETURN)

	if p.currToken != tokens.SEMI {
		p.parseAdditive()
	}

	p.expect(tokens.SEMI)
}

func (p *Parser) parseFor() {
	p.expect(tokens.FOR)
	p.expect(tokens.LPAREN)

	if p.currToken == tokens.VAR {
		p.parseVar()
	} else {
		p.parseAssignmentOrFuncCall(true)
	}

	p.parseComparison()
	p.expect(tokens.SEMI)
	p.parseAssignmentOrFuncCall(false)
	p.expect(tokens.RPAREN)
	p.parseBlock()
}

func (p *Parser) parseArguments() {
	for p.currToken != tokens.RPAREN {
		p.parseAdditive()

		if p.currToken == tokens.COMMA {
			p.advance()
		}
	}
}

func (p *Parser) parseComparison() {
	p.parseAdditive()

	for p.currToken == tokens.EQUAL || p.currToken == tokens.NEQUAL ||
		p.currToken == tokens.LT || p.currToken == tokens.LTOE ||
		p.currToken == tokens.GT || p.currToken == tokens.GTOE {
		p.advance()
		p.parseAdditive()
	}

	if p.currToken == tokens.ASSIGN {
		panic(fmt.Sprintf("unexpected assignment operator in comparison at %v", p.pos))
	}
}

func (p *Parser) parseAdditive() {
	p.parseMultiplicative()

	for p.currToken == tokens.ADD || p.currToken == tokens.SUB || p.currToken == tokens.DOT {
		p.advance()
		p.parseMultiplicative()
	}
}

func (p *Parser) parseMultiplicative() {
	p.parseFactor()

	for p.currToken == tokens.MUL || p.currToken == tokens.DIV || p.currToken == tokens.REM {
		p.advance()
		p.parseFactor()
	}
}

func (p *Parser) parseFactor() {
	switch p.currToken {
	case tokens.INT, tokens.STRING, tokens.FLOAT, tokens.CHAR, tokens.BOOL, tokens.TRUE, tokens.FALSE:
		p.advance()
	case tokens.IDENT:
		p.advance()

		switch p.currToken {
		case tokens.LPAREN:
			p.advance()
			p.parseArguments()
			p.expect(tokens.RPAREN)
		case tokens.ASSIGN:
			panic(fmt.Sprintf("unexpected assignment operator in factor at %v", p.pos))
		}
	default:
		panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
	}
}
