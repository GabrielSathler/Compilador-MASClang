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
	fmt.Println("position ", p.pos)
	fmt.Println("token ", p.currToken)
	fmt.Println("lexical ", p.currLex)
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

func (p *Parser) ParseProgram() {
	for p.currToken != tokens.EOF {
		switch p.currToken {
		case tokens.FUNC:
			fmt.Println("FUNC GLOBAL")
			p.parseFunction()
		case tokens.IF:
			fmt.Println("IF GLOBAL")
			p.parseIf()
		case tokens.VAR:
			fmt.Println("VAR GLOBAL")
			p.parseVar()
		case tokens.FOR:
			fmt.Println("FOR GLOBAL")
			p.parseFor()
		case tokens.WHILE:
			fmt.Println("WHILE GLOBAL")
			p.parseWhile()
		case tokens.PRINT:
			fmt.Println("PRINT GLOBAL")
			p.parsePrint()
		case tokens.INPUT:
			fmt.Println("INPUT GLOBAL")
			p.parseInput()
		case tokens.IDENT:
			fmt.Println("IDENT GLOBAL")
			p.parseStatement()
		default:
			panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
		}
	}
}

func (p *Parser) parseStatement() {
	fmt.Println(p.currToken)
	switch p.currToken {
	case tokens.IF:
		fmt.Println("IF")
		p.parseIf()
	case tokens.VAR:
		fmt.Println("VAR")
		p.parseVar()
	case tokens.FOR:
		fmt.Println("FOR")
		p.parseFor()
	case tokens.WHILE:
		fmt.Println("WHILE")
		p.parseWhile()
	case tokens.PRINT:
		fmt.Println("PRINT")
		p.parsePrint()
	case tokens.INPUT:
		fmt.Println("INPUT")
		p.parseInput()
	case tokens.IDENT:
		fmt.Println("IDENT")
		p.parseAssignmentOrFuncCall()
	case tokens.LBRACE:
		fmt.Println("LBRACE")
		p.parseBlock()
	case tokens.RETURN:
		fmt.Println("RETURN")
		p.parseReturn()
	default:
		panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
	}
}

func isValidType(t tokens.Token) bool {
	return t == tokens.INT || t == tokens.STRING || t == tokens.BOOL || t == tokens.CHAR || t == tokens.FLOAT
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
		panic(fmt.Sprintf("expected variable type, got %v at %v", p.currToken, p.pos))
	}

	p.advance()
	p.parseBlock()
}

func (p *Parser) parseFunctionParameters() {
	if p.currToken == tokens.RPAREN {
		return
	}

	for {
		if p.currToken != tokens.IDENT {
			panic(fmt.Sprintf("expected parameter name, got %v at %v", p.currToken, p.pos))
		}

		p.advance()
		p.expect(tokens.COLON)

		if !isValidType(p.currToken) {
			panic(fmt.Sprintf("expected parameter type, got %v at %v", p.currToken, p.pos))
		}

		p.advance()

		if p.currToken != tokens.COMMA {
			break
		}

		p.advance()
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
		p.parseComparison()
	}

	p.expect(tokens.SEMI)
}

func (p *Parser) parseFor() {
	p.expect(tokens.FOR)
	p.expect(tokens.LPAREN)

	if p.currToken == tokens.VAR {
		p.parseVar()
	} else {
		p.parseAssignmentOrFuncCall()
	}

	p.parseComparison()
	p.expect(tokens.SEMI)
	p.parseComparison()
	p.parseAssignmentOrFuncCall()
	p.expect(tokens.RPAREN)
	p.parseBlock()
}

func (p *Parser) parseAssignmentOrFuncCall() {
	p.advance()

	if p.currToken == tokens.ASSIGN {
		p.advance()
		p.parseComparison()
		p.expect(tokens.SEMI)
	} else if p.currToken == tokens.LPAREN {
		p.advance()
		p.parseArguments()
		p.expect(tokens.RPAREN)
		p.expect(tokens.SEMI)
	} else {
		panic(fmt.Sprintf("unexpected token after identifier %v at %v", p.currToken, p.pos))
	}
}

func (p *Parser) parseArguments() {
	if p.currToken == tokens.RPAREN {
		return
	}

	for {
		p.parseComparison()

		if p.currToken != tokens.COMMA {
			break
		}

		p.advance()
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
}

func (p *Parser) parseAdditive() {
	p.parseMultiplicative()

	for p.currToken == tokens.ADD || p.currToken == tokens.SUB {
		p.advance()
		p.parseMultiplicative()
	}
}

func (p *Parser) parseMultiplicative() {
	p.parseFactor()

	for p.currToken == tokens.MUL || p.currToken == tokens.DIV {
		p.advance()
		p.parseFactor()
	}
}

func (p *Parser) parseWhile() {
	p.expect(tokens.WHILE)
	p.expect(tokens.LPAREN)
	p.parseComparison()
	p.expect(tokens.RPAREN)
	p.parseBlock()
}

func (p *Parser) parsePrint() {
	p.expect(tokens.PRINT)
	p.expect(tokens.LPAREN)
	p.parseComparison()
	p.expect(tokens.RPAREN)
	p.expect(tokens.SEMI)
}

func (p *Parser) parseInput() {
	p.expect(tokens.INPUT)
	p.expect(tokens.LPAREN)

	if p.currToken != tokens.IDENT {
		panic(fmt.Sprintf("expected variable name in input, got %v at %v", p.currToken, p.pos))
	}

	p.advance()
	p.expect(tokens.RPAREN)
	p.expect(tokens.SEMI)
}

func (p *Parser) parseFactor() {
	switch p.currToken {
	case tokens.INT, tokens.CHAR, tokens.FLOAT, tokens.STRING, tokens.TRUE, tokens.FALSE:
		p.advance()
	case tokens.VAR:
		p.parseVar()
	case tokens.LPAREN:
		p.advance()
		p.parseComparison()
		p.expect(tokens.RPAREN)
	case tokens.IDENT:
		p.advance()
		if p.currToken == tokens.LPAREN {
			p.advance()
			p.parseArguments()
			p.expect(tokens.RPAREN)
		}
	default:
		panic(fmt.Sprintf("unexpected token %v at %v", p.currToken, p.pos))
	}
}

func (p *Parser) parseReturn() {
	p.expect(tokens.RETURN)

	if p.currToken != tokens.SEMI {
		p.parseComparison()
	}

	p.expect(tokens.SEMI)
}
