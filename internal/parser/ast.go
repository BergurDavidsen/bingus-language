package parser

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

type PrintStmt struct {
	Value Node
}

type LetStmt struct {
	Name  *IDent
	Value Node
}

type BoolLit struct {
	Value bool
}

type BinaryExpr struct {
	Left     Node
	Operator string
	Right    Node
}

type UnaryExpr struct {
	Operator string
	Right    Node
}
