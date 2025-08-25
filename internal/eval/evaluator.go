package eval

import (
	"fmt"
	"strconv"

	"github.com/BergurDavidsen/bingus/internal/parser"
)

type Env struct {
	vars map[string]int
}

func NewEnv() *Env {
	return &Env{vars: make(map[string]int)}
}

func (e *Env) Eval(node parser.Node) int {
	switch n := node.(type) {
	case *parser.Program:
		var result int
		for _, stmt := range n.Statements {
			result = e.Eval(stmt)
		}
		return result

	case *parser.ReturnStmt:
		return e.Eval(n.Value)

	case *parser.NumberLiteral:
		val, _ := strconv.Atoi(n.Value)
		return val

	case *parser.IDent:
		return e.vars[n.Name]

	case *parser.LetStmt:
		val := e.Eval(n.Value)
		e.vars[n.Name.Name] = val
		return val

	case *parser.BinaryExpr:
		left := e.Eval(n.Left)
		right := e.Eval(n.Right)
		switch n.Operator {
		case "+":
			return left + right
		case "-":
			return left - right
		case "/":
			return left / right
		case "*":
			return left * right
		default:
			panic("unknown operator " + n.Operator)
		}

	case *parser.UnaryExpr:
		right := e.Eval(n.Right)

		switch n.Operator {
		case "+":
			return +right
		case "-":
			return -right

		default:
			panic("unknown unary operator " + n.Operator)
		}

	default:
		panic(fmt.Sprintf("unhandled node type: %T", n))
	}
}
