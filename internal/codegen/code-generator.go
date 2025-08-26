package codegen

import (
	"fmt"
	"strings"

	"github.com/BergurDavidsen/bingus/internal/parser"
)

type CodeGen struct {
	code     []string
	symbols  map[string]int
	stackPos int
}

func NewCodeGen() *CodeGen {
	return &CodeGen{
		code:     []string{},
		symbols:  make(map[string]int),
		stackPos: 0,
	}
}

func (cg *CodeGen) Emit(line string) {
	cg.code = append(cg.code, line)
}

func (cg *CodeGen) EmitIndent(indent int, line string) {
	cg.code = append(cg.code, strings.Repeat("  ", indent)+line)
}

func (cg *CodeGen) String() string {
	return strings.Join(cg.code, "\n")
}

func (cg *CodeGen) Gen(node parser.Node) {
	// Program prologue
	cg.Emit("section .text")
	cg.Emit("global _start")
	cg.Emit("_start:")

	cg.EmitIndent(1, "push rbp")
	cg.EmitIndent(1, "mov rbp, rsp")

	switch n := node.(type) {
	case *parser.Program:
		for _, stmt := range n.Statements {
			cg.GenStmt(stmt)
		}
	}

}

func (cg *CodeGen) GenStmt(node parser.Node) {
	switch n := node.(type) {
	case *parser.ReturnStmt:
		val := cg.GenExpr(n.Value)
		cg.EmitIndent(1, fmt.Sprintf("mov rdi, %s", val)) // exit code
		cg.EmitIndent(1, "mov rax, 60")                   // syscall: exit
		cg.EmitIndent(1, "syscall")
	case *parser.LetStmt:
		val := cg.GenExpr(n.Value)

		cg.stackPos += 8
		offset := cg.stackPos
		cg.symbols[n.Name.Name] = offset

		cg.EmitIndent(1, fmt.Sprintf("mov QWORD [rbp-%d], %s", offset, val))

	default:
		panic(fmt.Sprintf("unsupported statement: %T", n))
	}
}
func (cg *CodeGen) GenExpr(node parser.Node) string {
	return cg.genExprWithTarget(node, "rax")
}

func (cg *CodeGen) genExprWithTarget(node parser.Node, target string) string {
	switch n := node.(type) {
	case *parser.NumberLiteral:
		cg.EmitIndent(1, fmt.Sprintf("mov %s, %s", target, n.Value))
		return n.Value
	case *parser.IDent:
		offset, ok := cg.symbols[n.Name]
		if !ok {
			panic(fmt.Sprintf("undefined variable: %s", n.Name))
		}

		cg.EmitIndent(1, fmt.Sprintf("mov %s, [rbp-%d]", target, offset))
		return target

	case *parser.UnaryExpr:
		cg.GenExpr(n.Right)
		op := n.Operator
		if op == "-" {
			cg.EmitIndent(1, fmt.Sprintf("neg %s", target))
		}
		return target

	case *parser.BinaryExpr:

		op := n.Operator
		if op == "/" {
			cg.GenExpr(n.Left)
			cg.EmitIndent(1, "xor rdx, rdx")
			cg.genExprWithTarget(n.Right, "rbx")
			cg.EmitIndent(1, "idiv rbx")
			return target

		}
		cg.GenExpr(n.Left)

		cg.EmitIndent(1, fmt.Sprintf("push %s", target))
		cg.GenExpr(n.Right)
		cg.EmitIndent(1, "pop rbx")

		switch op {
		case "+":
			cg.EmitIndent(1, fmt.Sprintf("add %s, rbx", target))
		case "-":
			cg.EmitIndent(1, fmt.Sprintf("sub rbx, %s", target))
			cg.EmitIndent(1, fmt.Sprintf("mov %s, rbx", target))
		case "*":
			cg.EmitIndent(1, fmt.Sprintf("mul %s, rbx", target))
		}
		return target

	default:
		panic(fmt.Sprintf("unsupported expression: %T", n))
	}

}
