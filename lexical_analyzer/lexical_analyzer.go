package lexical_analyzer

import (
	"bufio"
	"io"
	"unicode"

	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, tokens.Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, tokens.EOF, ""
			}
			panic(err)
		}

		l.pos.Column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, tokens.SEMI, ";"
		case '(':
			return l.pos, tokens.LPAREN, "("
		case ')':
			return l.pos, tokens.RPAREN, ")"
		case '{':
			return l.pos, tokens.LBRACE, "{"
		case '}':
			return l.pos, tokens.RBRACE, "}"
		case ':':
			return l.pos, tokens.COLON, ":"
		case ',':
			return l.pos, tokens.COMMA, ","
		case '+':
			return l.pos, tokens.ADD, "+"
		case '-':
			return l.pos, tokens.SUB, "-"
		case '*':
			return l.pos, tokens.MUL, "*"
		case '/':
			return l.pos, tokens.DIV, "/"
		case '=':
			return l.pos, tokens.ASSIGN, "="
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, tokens.INT, lit
			} else if unicode.IsLetter(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, tokens.IDENT, lit
			} else {
				return l.pos, tokens.ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.Line++
	l.pos.Column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.Column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.Column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}
