package main

type Node interface{}

type Program struct {
	Statements []Node
}

type ReturnStmt struct {
	Value Node
}

type NumberLiteral struct {
	Value string
}

type IDent struct {
	Name string
}

type LetStmt struct {
	Name  *IDent
	Value Node
}

type BinaryExpr struct {
	Left     Node
	Operator string
	Right    Node
}
