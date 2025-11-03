package lexer

const (
	TOKEN_IF = iota
	TOKEN_ELSE
	TOKEN_WHILE
	TOKEN_FOR
	TOKEN_RETURN
	TOKEN_LET
	TOKEN_IDENT
	TOKEN_NUMBER
	TOKEN_PRINT
	TOKEN_STRING
	TOKEN_EQUAL
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_SEMICOLON
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULTIPLY
	TOKEN_DIVIDE
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_LT
	TOKEN_GT
	TOKEN_LE
	TOKEN_GE
	TOKEN_EQ
)

var keywords = map[string]int{
	"if":     TOKEN_IF,
	"else":   TOKEN_ELSE,
	"while":  TOKEN_WHILE,
	"for":    TOKEN_FOR,
	"let":    TOKEN_LET,
	"true":   TOKEN_TRUE,
	"false":  TOKEN_FALSE,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
}

var multiCharTokens = map[string]int{
	"==": TOKEN_EQ,
	"<=": TOKEN_LE,
	">=": TOKEN_GE,
}

var singleCharTokens = map[byte]int{
	'+': TOKEN_PLUS,
	'-': TOKEN_MINUS,
	'*': TOKEN_MULTIPLY,
	'/': TOKEN_DIVIDE,
	'=': TOKEN_EQUAL,
	';': TOKEN_SEMICOLON,
	'(': TOKEN_LPAREN,
	')': TOKEN_RPAREN,
	'{': TOKEN_LBRACE,
	'}': TOKEN_RBRACE,
	'<': TOKEN_LT,
	'>': TOKEN_GT,
}

type Token struct {
	Type    int
	Literal string
}
