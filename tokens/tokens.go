package tokens

type Token int

const (
	EOF Token = iota
	ILLEGAL

	IDENT
	INT
	FLOAT
	CHAR
	BOOL
	STRING

	ASSIGN
	ADD
	SUB
	MUL
	DIV

	SEMI
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	COLON
	COMMA

	VAR
	FUNC
	IF
	ELSE
	WHILE
	FOR
	RETURN
	PRINT
	INPUT
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	BOOL:   "BOOL",
	STRING: "STRING",

	ASSIGN: "=",
	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",

	SEMI:   ";",
	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",
	COLON:  ":",
	COMMA:  ",",

	VAR:    "var",
	FUNC:   "func",
	IF:     "if",
	ELSE:   "else",
	WHILE:  "while",
	FOR:    "for",
	RETURN: "return",
	PRINT:  "print",
	INPUT:  "input",
}

func (t Token) String() string {
	return tokens[t]
}
