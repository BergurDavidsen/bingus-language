package main

import (
	"fmt"
	"unicode"
)

func lex(input string) []Token {
	var tokens []Token

	i := 0

	for i < len(input) {
		c := input[i]

		// Skip whitespace
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
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
