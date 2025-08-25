package codegen

import (
	"fmt"
	"strings"

	"github.com/BergurDavidsen/bingus/internal/parser"
)

type CodeGen struct {
	code []string
}

func NewCodeGen() *CodeGen {
	return &CodeGen{code: []string{}}
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
	default:
		panic(fmt.Sprintf("unsupported statement: %T", n))
	}
}

func (cg *CodeGen) GenExpr(node parser.Node) string {
	switch n := node.(type) {
	case *parser.NumberLiteral:
		return n.Value
	default:
		panic(fmt.Sprintf("unsupported expression: %T", n))
	}

}
