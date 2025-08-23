package main

import (
	"fmt"
	"reflect"
)

type Parser struct {
	tokens []Token
	pos    int
}

func printNodeReflect(node interface{}, indent string) {
	if node == nil {
		fmt.Println(indent + "nil")
		return
	}

	v := reflect.ValueOf(node)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Println(indent + "nil")
			return
		}
		printNodeReflect(v.Elem().Interface(), indent)
		return
	}

	switch v.Kind() {
	case reflect.Struct:
		fmt.Println(indent + v.Type().Name())
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			val := v.Field(i).Interface()
			fmt.Printf("%s  %s: ", indent, field.Name)
			kind := reflect.ValueOf(val).Kind()
			if kind == reflect.Struct || kind == reflect.Ptr || kind == reflect.Slice {
				fmt.Println()
				printNodeReflect(val, indent+"    ")
			} else {
				fmt.Printf("%v\n", val)
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			printNodeReflect(v.Index(i).Interface(), indent+"  ")
		}

	default:
		fmt.Println(indent, v.Interface())
	}
}

func (p *Parser) currentToken() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: -1, Literal: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) parseNumber() *NumberLiteral {
	tok := p.currentToken()

	if tok.Type != TOKEN_NUMBER {
		panic(fmt.Sprintf("Expected number, got: %v", tok))
	}

	p.advance()
	return &NumberLiteral{Value: tok.Literal}
}

func (p *Parser) parseIdent() *IDent {
	tok := p.currentToken()

	if tok.Type != TOKEN_IDENT {
		panic("Expected identifier after let")
	}
	p.advance()

	return &IDent{Name: tok.Literal}
}
func (p *Parser) parseReturnStmt() *ReturnStmt {
	tok := p.currentToken()
	if tok.Type != TOKEN_RETURN {
		panic(fmt.Sprintf("Expected 'return', got: %v", tok))
	}
	p.advance()

	value := p.parserExpression()

	if p.currentToken().Type != TOKEN_SEMICOLON {
		panic(fmt.Sprintf("Expected ';', got: %v", p.currentToken()))
	}
	p.advance()

	return &ReturnStmt{Value: value}
}

func (p *Parser) parseLetStmt() *LetStmt {
	tok := p.currentToken()

	if tok.Type != TOKEN_LET {
		panic(fmt.Sprintf("Expected 'let', got: %v", tok))
	}

	p.advance() // ignore let

	if p.currentToken().Type != TOKEN_IDENT {
		panic("Expected identifier after let")
	}

	id := p.parseIdent()

	if p.currentToken().Type != TOKEN_EQUAL {
		panic("Expected '=' after identifier in let statement")
	}

	p.advance()

	value := p.parserExpression()

	if p.currentToken().Type != TOKEN_SEMICOLON {
		panic(fmt.Sprintf("Expected ';', got: %v", tok))
	}
	p.advance()

	return &LetStmt{Name: id, Value: value}
}

func (p *Parser) parserExpression() Node {

	left := p.parsePrimary()

	for p.pos < len(p.tokens) {
		tok := p.currentToken()
		if tok.Literal != "+" && tok.Literal != "-" {
			break
		}
		op := tok.Literal
		p.advance()

		right := p.parsePrimary()

		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}
	return left

}

func (p *Parser) parsePrimary() Node {
	tok := p.currentToken()

	switch tok.Type {
	case TOKEN_NUMBER:
		return p.parseNumber()
	case TOKEN_IDENT:
		return p.parseIdent()
	default:
		panic(fmt.Sprintf("Expected number after return, got: %v", p.currentToken()))
	}
}

func (p *Parser) parseProgram() *Program {
	prog := &Program{}

	for p.pos < len(p.tokens) {
		tok := p.currentToken()
		switch tok.Type {
		case TOKEN_RETURN:
			stmt := p.parseReturnStmt()
			prog.Statements = append(prog.Statements, stmt)
		case TOKEN_LET:
			stmt := p.parseLetStmt()
			prog.Statements = append(prog.Statements, stmt)
		default:
			panic(fmt.Sprintf("Unexpected token: %v", tok))
		}
	}
	return prog
}
