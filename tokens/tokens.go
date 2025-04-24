package tokens

type Token int

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	CHAR   = "CHAR"
	BOOL   = "BOOL"
	STRING = "STRING"

	ASSIGN = "="
	ADD    = "+"
	SUB    = "-"
	MUL    = "*"
	DIV    = "/"

	SEMI   = ";"
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	COLON  = ":"
	COMMA  = ","

	VAR    = "var"
	FUNC   = "func"
	IF     = "if"
	ELSE   = "else"
	WHILE  = "while"
	FOR    = "for"
	RETURN = "return"
	PRINT  = "print"
	INPUT  = "input"
)

var Tokens = map[string]string{
	"var":    VAR,
	"func":   FUNC,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"return": RETURN,
	"print":  PRINT,
	"input":  INPUT,

	"int":    INT,
	"float":  FLOAT,
	"char":   CHAR,
	"bool":   BOOL,
	"string": STRING,

	"true":  BOOL,
	"false": BOOL,
}
