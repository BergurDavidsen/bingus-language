package parser

import (
	"fmt"
	"reflect"

	"github.com/BergurDavidsen/bingus/internal/lexer"
)

type Parser struct {
	Tokens []lexer.Token
	pos    int
}

var precedences = map[string]int{
	"==": 1, // lowest precedence
	"<":  1,
	"<=": 1,
	">":  1,
	">=": 1,
	"+":  2,
	"-":  2,
	"*":  3,
	"/":  3,
	"u+": 4,
	"u-": 4,
}

func getPrecedence(tok lexer.Token) int {
	if prec, ok := precedences[tok.Literal]; ok {
		return prec
	}
	return 0
}

func PrintNodeReflect(node interface{}, indent string) {
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
		PrintNodeReflect(v.Elem().Interface(), indent)
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
				PrintNodeReflect(val, indent+"    ")
			} else {
				fmt.Printf("%v\n", val)
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			PrintNodeReflect(v.Index(i).Interface(), indent+"  ")
		}

	default:
		fmt.Println(indent, v.Interface())
	}
}

func (p *Parser) currentToken() lexer.Token {
	if p.pos >= len(p.Tokens) {
		return lexer.Token{Type: -1, Literal: ""}
	}
	return p.Tokens[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) parseNumber() *NumberLiteral {
	tok := p.currentToken()

	if tok.Type != lexer.TOKEN_NUMBER {
		panic(fmt.Sprintf("Expected number, got: %v", tok))
	}

	p.advance()
	return &NumberLiteral{Value: tok.Literal}
}

func (p *Parser) parseIdent() *IDent {
	tok := p.currentToken()

	if tok.Type != lexer.TOKEN_IDENT {
		panic("Expected identifier after let")
	}
	p.advance()

	return &IDent{Name: tok.Literal}
}
func (p *Parser) parseReturnStmt() *ReturnStmt {
	tok := p.currentToken()
	if tok.Type != lexer.TOKEN_RETURN {
		panic(fmt.Sprintf("Expected 'return', got: %v", tok))
	}
	p.advance()

	value := p.parserExpression(1)

	if p.currentToken().Type != lexer.TOKEN_SEMICOLON {
		panic(fmt.Sprintf("Expected ';', got: %v", p.currentToken()))
	}
	p.advance()

	return &ReturnStmt{Value: value}
}

func (p *Parser) parsePrint() *PrintStmt {
	tok := p.currentToken()

	if tok.Type != lexer.TOKEN_PRINT {
		panic(fmt.Sprintf("Expected 'print', got: %v", tok))
	}
	p.advance()

	value := p.parserExpression(1)

	if p.currentToken().Type != lexer.TOKEN_SEMICOLON {
		panic(fmt.Sprintf("Expected ';', got: %v", p.currentToken()))
	}
	p.advance()

	return &PrintStmt{Value: value}
}

func (p *Parser) parseLetStmt() *LetStmt {
	tok := p.currentToken()

	if tok.Type != lexer.TOKEN_LET {
		panic(fmt.Sprintf("Expected 'let', got: %v", tok))
	}

	p.advance() // ignore let

	if p.currentToken().Type != lexer.TOKEN_IDENT {
		panic("Expected identifier after let")
	}

	id := p.parseIdent()

	if p.currentToken().Type != lexer.TOKEN_EQUAL {
		panic("Expected '=' after identifier in let statement")
	}

	p.advance()

	value := p.parserExpression(1)

	if p.currentToken().Type != lexer.TOKEN_SEMICOLON {
		panic(fmt.Sprintf("Expected ';', got: %v", tok))
	}
	p.advance()

	return &LetStmt{Name: id, Value: value}
}
func (p *Parser) parseBlock() []Node {
	stmts := []Node{}

	if p.currentToken().Type != lexer.TOKEN_LBRACE {
		panic("Expected '{' at start of block")
	}
	p.advance() // consume '{'

	for p.currentToken().Type != lexer.TOKEN_RBRACE {
		tok := p.currentToken()
		switch tok.Type {
		case lexer.TOKEN_RETURN:
			stmts = append(stmts, p.parseReturnStmt())
		case lexer.TOKEN_LET:
			stmts = append(stmts, p.parseLetStmt())
		case lexer.TOKEN_PRINT:
			stmts = append(stmts, p.parsePrint())
		case lexer.TOKEN_IF:
			stmts = append(stmts, p.parseIfStmt())
		default:
			panic(fmt.Sprintf("Unexpected token in block: %v", tok))
		}
	}

	p.advance() // consume '}'
	return stmts
}

func (p *Parser) parseIfStmt() *IfStmt {
	tok := p.currentToken()

	if tok.Type != lexer.TOKEN_IF {
		panic(fmt.Sprintf("Expected 'if', got: %v", tok))
	}

	p.advance()

	if p.currentToken().Type != lexer.TOKEN_LPAREN {
		panic("Expected '(' after if statement")
	}

	p.advance()

	gaurd := p.parserExpression(1)

	if p.currentToken().Type != lexer.TOKEN_RPAREN {
		panic("Expected ')' after if statement")
	}
	p.advance()

	thenBlock := p.parseBlock()

	if p.currentToken().Type != lexer.TOKEN_ELSE {
		panic("Expected 'else' after if statement")
	}
	p.advance()

	elseBlock := p.parseBlock()

	return &IfStmt{Guard: gaurd, Then: thenBlock, Else: elseBlock}
}

func (p *Parser) parserExpression(minPrec int) Node {

	left := p.parsePrimary()

	for p.pos < len(p.Tokens) {
		tok := p.currentToken()
		prec := getPrecedence(tok)

		if prec < minPrec {
			break
		}

		op := tok.Literal
		p.advance()

		right := p.parserExpression(prec)

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
	case lexer.TOKEN_NUMBER:
		return p.parseNumber()
	case lexer.TOKEN_IDENT:
		return p.parseIdent()
	case lexer.TOKEN_LPAREN:
		p.advance()
		expr := p.parserExpression(1)
		if p.currentToken().Type != lexer.TOKEN_RPAREN {
			panic(fmt.Sprintf("Expected ')', got: %v", p.currentToken()))
		}
		p.advance()
		return expr
	case lexer.TOKEN_PLUS, lexer.TOKEN_MINUS:
		op := tok.Literal
		p.advance()
		right := p.parsePrimary()
		return &UnaryExpr{
			Operator: op,
			Right:    right,
		}
	case lexer.TOKEN_TRUE, lexer.TOKEN_FALSE:
		val := tok.Type == lexer.TOKEN_TRUE
		p.advance()
		return &BoolLit{Value: val}
	default:
		panic(fmt.Sprintf("Expected number after return, got: %v", p.currentToken()))
	}
}

func (p *Parser) ParseProgram() *Program {
	prog := &Program{}

	for p.pos < len(p.Tokens) {
		tok := p.currentToken()
		switch tok.Type {
		case lexer.TOKEN_RETURN:
			stmt := p.parseReturnStmt()
			prog.Statements = append(prog.Statements, stmt)
		case lexer.TOKEN_LET:
			stmt := p.parseLetStmt()
			prog.Statements = append(prog.Statements, stmt)
		case lexer.TOKEN_PRINT:
			stmt := p.parsePrint()
			prog.Statements = append(prog.Statements, stmt)
		case lexer.TOKEN_IF:
			stmt := p.parseIfStmt()
			prog.Statements = append(prog.Statements, stmt)
		default:
			panic(fmt.Sprintf("Unexpected lexer.token: %v", tok))
		}
	}
	return prog
}
