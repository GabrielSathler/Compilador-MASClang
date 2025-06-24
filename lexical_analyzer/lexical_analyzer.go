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
		currentRune, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, tokens.EOF, ""
			}

			panic(err)
		}

		l.pos.Column++

		switch currentRune {
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
			next, _, err := l.reader.ReadRune()

			if err == nil && next == '=' {
				l.pos.Column++
				return l.pos, tokens.EQUAL, "=="
			}

			l.backup()
			return l.pos, tokens.ASSIGN, "="
		case '!':
			next, _, err := l.reader.ReadRune()

			if err == nil && next == '=' {
				l.pos.Column++
				return l.pos, tokens.NEQUAL, "!="
			}

			l.backup()
			return l.pos, tokens.NOT, "!"
		case '<':
			next, _, err := l.reader.ReadRune()

			if err == nil && next == '=' {
				l.pos.Column++
				return l.pos, tokens.LTOE, "<="
			}

			l.backup()
			return l.pos, tokens.LT, "<"
		case '>':
			next, _, err := l.reader.ReadRune()

			if err == nil && next == '=' {
				l.pos.Column++
				return l.pos, tokens.GTOE, ">="
			}

			l.backup()
			return l.pos, tokens.GT, ">"
		case '"':
			startPos := l.pos
			lit := l.lexString()
			return startPos, tokens.STRING, lit
		case '\'':
			startPos := l.pos
			lit := l.lexChar()
			return startPos, tokens.CHAR, lit
		default:
			if currentRune == '_' {
				startPos := l.pos
				lit := l.lexIdentWithUnderscore()

				return startPos, tokens.IDENT, lit
			}

			if unicode.IsSpace(currentRune) {
				continue
			} else if unicode.IsDigit(currentRune) {
				startPos := l.pos
				l.backup()
				lit, tokenType := l.lexNumber()

				return startPos, tokenType, lit
			} else if unicode.IsLetter(currentRune) {
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()

				switch lit {
				case "func":
					return startPos, tokens.FUNC, lit
				case "var":
					return startPos, tokens.VAR, lit
				case "return":
					return startPos, tokens.RETURN, lit
				case "int":
					return startPos, tokens.INT, lit
				case "float":
					return startPos, tokens.FLOAT, lit
				case "char":
					return startPos, tokens.CHAR, lit
				case "bool":
					return startPos, tokens.BOOL, lit
				case "string":
					return startPos, tokens.STRING, lit
				case "for":
					return startPos, tokens.FOR, lit
				case "while":
					return startPos, tokens.WHILE, lit
				case "if":
					return startPos, tokens.IF, lit
				case "else":
					return startPos, tokens.ELSE, lit
				case "true":
					return startPos, tokens.TRUE, lit
				case "false":
					return startPos, tokens.FALSE, lit
				case "print":
					return startPos, tokens.PRINT, lit
				case "input":
					return startPos, tokens.INPUT, lit
				default:
					return startPos, tokens.IDENT, lit
				}
			} else if currentRune == '%' {
				return l.pos, tokens.REM, "%"
			} else if currentRune == '.' {
				return l.pos, tokens.DOT, "."
			} else {
				return l.pos, tokens.ILLEGAL, string(currentRune)
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

func (l *Lexer) lexNumber() (string, tokens.Token) {
	var lit string
	isFloat := false

	for {
		currentRune, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit, tokens.INT
			}
		}

		if unicode.IsDigit(currentRune) {
			lit += string(currentRune)
		} else if currentRune == '.' && !isFloat {
			isFloat = true
			lit += string(currentRune)

			nextRune, _, err := l.reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return lit, tokens.ILLEGAL
				}
			}

			if !unicode.IsDigit(nextRune) {
				l.backup()
				return lit, tokens.ILLEGAL
			}

			lit += string(nextRune)
		} else {
			l.backup()

			if isFloat {
				return lit, tokens.FLOAT
			}

			return lit, tokens.INT
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string

	for {
		currentRune, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.Column++

		if unicode.IsLetter(currentRune) || unicode.IsDigit(currentRune) || currentRune == '_' {
			lit = lit + string(currentRune)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdentWithUnderscore() string {
	var lit string

	for {
		currentRune, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.Column++

		if unicode.IsLetter(currentRune) || unicode.IsDigit(currentRune) || currentRune == '_' {
			lit = lit + string(currentRune)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexString() string {
	var lit string

	for {
		currentRune, _, err := l.reader.ReadRune()
		if err != nil {
			panic("unterminated string literal")
		}

		l.pos.Column++

		if currentRune == '"' {
			break
		}

		lit += string(currentRune)
	}
	return lit
}

func (l *Lexer) lexChar() string {
	currentRune, _, err := l.reader.ReadRune()
	if err != nil {
		panic("unterminated char literal")
	}

	l.pos.Column++

	if currentRune == '\'' {
		panic("empty char literal")
	}

	lit := string(currentRune)
	nextRune, _, err := l.reader.ReadRune()

	if err != nil || nextRune != '\'' {
		panic("unterminated or invalid char literal")
	}

	l.pos.Column++

	return lit
}
