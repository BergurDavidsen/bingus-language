package lexer

import (
	"fmt"
	"unicode"
)

func Lex(input string) []Token {
	var tokens []Token

	i := 0

	for i < len(input) {
		c := input[i]

		// Skip whitespace
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
		}

		if c == '/' && i+1 < len(input) {
			next := input[i+1]

			// Single-line comment //
			if next == '/' {
				i += 2
				for i < len(input) && input[i] != '\n' {
					i++
				}
				continue
			}

			// Multi-line comment /* ... */
			if next == '*' {
				i += 2
				for i+1 < len(input) && !(input[i] == '*' && input[i+1] == '/') {
					i++
				}
				if i+1 >= len(input) {
					panic("Unterminated multi-line comment")
				}
				i += 2
				continue
			}
		}

		if i+1 < len(input) {
			twoChar := input[i : i+2]
			if tokType, ok := multiCharTokens[twoChar]; ok {
				tokens = append(tokens, Token{Type: tokType, Literal: twoChar})
				i += 2
				continue
			}
		}

		if tokType, ok := singleCharTokens[c]; ok {
			tokens = append(tokens, Token{Type: tokType, Literal: string(c)})
			i++
			continue
		}
		if c == '"' {
			j := i + 1
			for j < len(input) && input[j] != '"' {
				j++
			}
			if j >= len(input) {
				panic("Unterminated string literal")
			}

			str := input[i+1 : j]
			tokens = append(tokens, Token{Type: TOKEN_STRING, Literal: str})
			i = j + 1
			continue
		}

		if unicode.IsDigit(rune(c)) {
			j := i
			for j < len(input) && unicode.IsDigit(rune(input[j])) {
				j++
			}
			num := input[i:j]
			tokens = append(tokens, Token{Type: TOKEN_NUMBER, Literal: num})
			i = j
			continue
		}

		if unicode.IsLetter(rune(c)) {
			j := i
			for j < len(input) && (unicode.IsLetter(rune(input[j])) || unicode.IsDigit(rune(input[j])) || unicode.IsLetter('_')) {
				j++
			}
			word := input[i:j]

			if tokType, ok := keywords[word]; ok {
				tokens = append(tokens, Token{Type: tokType, Literal: word})
			} else {
				tokens = append(tokens, Token{Type: TOKEN_IDENT, Literal: word})
			}

			i = j
			continue
		}

		panic(fmt.Sprintf("Unexpected character: %c", c))
	}
	return tokens
}
