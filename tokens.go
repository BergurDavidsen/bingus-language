package main

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
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_LBRACE    // {
	TOKEN_RBRACE    // }
	TOKEN_SEMICOLON // ;
	TOKEN_PLUS
	TOKEN_MINUS
)

var keywords = map[string]int{
	"if":     TOKEN_IF,
	"else":   TOKEN_ELSE,
	"while":  TOKEN_WHILE,
	"for":    TOKEN_FOR,
	"let":    TOKEN_LET,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
}

var singleCharTokens = map[byte]int{
	'+': TOKEN_PLUS,
	'-': TOKEN_MINUS,
	'=': TOKEN_EQUAL,
	';': TOKEN_SEMICOLON,
	'(': TOKEN_LPAREN,
	')': TOKEN_RPAREN,
}

type Token struct {
	Type    int
	Literal string
}
