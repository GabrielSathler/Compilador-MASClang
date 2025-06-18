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
	TRUE
	FALSE

	ASSIGN
	ADD
	SUB
	MUL
	DIV

	EQUAL
	NEQUAL
	LT
	LTOE
	GT
	GTOE
	AND
	OR

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
	INT:    "int",
	FLOAT:  "float",
	CHAR:   "char",
	BOOL:   "bool",
	STRING: "string",
	TRUE:   "true",
	FALSE:  "false",

	ASSIGN: "=",
	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",
	EQUAL:  "==",
	NEQUAL: "!=",
	LT:     "<",
	LTOE:   "<=",
	GT:     ">",
	GTOE:   ">=",
	AND:    "&&",
	OR:     "||",

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
