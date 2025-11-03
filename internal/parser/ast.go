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

type AssignmentStmt struct {
	Name  *IDent
	Value Node
}

type PrintStmt struct {
	Value Node
}

type LetStmt struct {
	Name  *IDent
	Value Node
}

type WhileStmt struct {
	Guard Node
	Body  []Node
}

type BoolLit struct {
	Value bool
}

type IfStmt struct {
	Guard Node
	Then  []Node
	Else  []Node
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

type BreakStmt struct{}

type ContinueStmt struct{}
