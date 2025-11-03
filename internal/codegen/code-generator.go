package codegen

import (
	"fmt"
	"strings"

	"github.com/BergurDavidsen/bingus/internal/parser"
)

type CodeGen struct {
	code           []string
	scope          []map[string]int
	stackMark      []int
	stackPos       int
	labelCnt       int
	loopStartStack []string
	loopEndStack   []string
}

func (cg *CodeGen) newLabel(base string) string {
	cg.labelCnt++
	return fmt.Sprintf(".%s_%d", base, cg.labelCnt)
}

func NewCodeGen() *CodeGen {
	return &CodeGen{
		code:      []string{},
		scope:     []map[string]int{{}},
		stackMark: []int{},
		stackPos:  0,
		labelCnt:  0,
	}
}

func (cg *CodeGen) pushScope() {
	cg.scope = append(cg.scope, map[string]int{})
	cg.stackMark = append(cg.stackMark, cg.stackPos)
}

func (cg *CodeGen) popScope() {
	if len(cg.scope) == 0 {
		panic("no scope to pop")
	}

	prevMark := cg.stackMark[len(cg.stackMark)-1] // saved at push
	// compute how much space was allocated while in this scope
	delta := cg.stackPos - prevMark
	if delta > 0 {
		cg.EmitIndent(1, fmt.Sprintf("add rsp, %d", delta))
	}

	cg.scope = cg.scope[:len(cg.scope)-1]
	cg.stackMark = cg.stackMark[:len(cg.stackMark)-1]
	cg.stackPos = prevMark
}

func (cg *CodeGen) currentScope() map[string]int {
	return cg.scope[len(cg.scope)-1]
}

func (cg *CodeGen) declareVar(name string, offset int) {
	scope := cg.currentScope()

	if _, exists := scope[name]; exists {
		panic(fmt.Sprintf("variable already declared in this scope: %s", name))
	}

	scope[name] = offset
}

func (cg *CodeGen) lookupVar(name string) (int, bool) {
	for i := len(cg.scope) - 1; i >= 0; i-- {
		if offset, ok := cg.scope[i][name]; ok {
			return offset, true
		}
	}
	return 0, false
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
		lastStmt := n.Statements[len(n.Statements)-1]

		if _, ok := lastStmt.(*parser.ReturnStmt); !ok {
			fmt.Println("Warning: no return statement at end of program; adding a default 'return 0' to end of file.")
			n.Statements = append(n.Statements, &parser.ReturnStmt{
				Value: &parser.NumberLiteral{Value: "0"},
			})
		}
		for _, stmt := range n.Statements {
			cg.GenStmt(stmt)
		}
	}

	asmHelper := `
		section .bss
		buffer resb 20

		section .text
		print_number:
			mov rax, rdi
			mov rcx, 0
			lea rbx, [buffer+19]
			mov byte [rbx], 10
			dec rbx

		.convert_loop:
			xor rdx, rdx
			mov rsi, 10
			div rsi
			add dl, '0'
			mov [rbx], dl
			dec rbx
			inc rcx
			test rax, rax
			jnz .convert_loop

			inc rbx

			mov rax, 1
			mov rdi, 1
			mov rsi, rbx
			inc rcx
			mov rdx, rcx
			syscall

			ret
		`

	cg.Emit(asmHelper)

}

func (cg *CodeGen) GenStmt(node parser.Node) {
	switch n := node.(type) {
	case *parser.ReturnStmt:
		val := cg.GenExpr(n.Value)
		cg.EmitIndent(1, fmt.Sprintf("mov rdi, %s", val)) // exit code

		// Tear down stack frame before exit
		cg.EmitIndent(1, "mov rsp, rbp")
		cg.EmitIndent(1, "pop rbp")

		cg.EmitIndent(1, "mov rax, 60") // syscall: exit
		cg.EmitIndent(1, "syscall")
	case *parser.LetStmt:
		val := cg.GenExpr(n.Value)

		cg.stackPos += 8
		cg.EmitIndent(1, "sub rsp, 8")
		offset := cg.stackPos
		cg.declareVar(n.Name.Name, offset)

		cg.EmitIndent(1, fmt.Sprintf("mov QWORD [rbp-%d], %s", offset, val))

	case *parser.PrintStmt:
		val := cg.GenExpr(n.Value)
		cg.EmitIndent(1, fmt.Sprintf("mov rdi, %s", val))
		cg.EmitIndent(1, "call print_number")

	case *parser.IfStmt:
		cg.GenExpr(n.Guard)

		elseLabel := cg.newLabel("else")
		end_label := cg.newLabel("endif")

		cg.EmitIndent(1, "cmp rax, 0")
		cg.EmitIndent(1, fmt.Sprintf("je %s", elseLabel))

		cg.pushScope()
		for _, stmt := range n.Then {
			cg.GenStmt(stmt)
		}
		cg.popScope()

		cg.EmitIndent(1, fmt.Sprintf("jmp %s", end_label))

		cg.Emit(fmt.Sprintf("%s:", elseLabel))

		cg.pushScope()
		for _, stmt := range n.Else {
			cg.GenStmt(stmt)
		}
		cg.popScope()

		cg.Emit(fmt.Sprintf("%s:", end_label))
	case *parser.WhileStmt:
		start_label := cg.newLabel("while_start")
		end_label := cg.newLabel("while_end")

		cg.loopStartStack = append(cg.loopStartStack, start_label)
		cg.loopEndStack = append(cg.loopEndStack, end_label)

		cg.Emit(fmt.Sprintf("%s:", start_label))

		cg.GenExpr(n.Guard)

		cg.EmitIndent(1, "cmp rax, 0")
		cg.EmitIndent(1, fmt.Sprintf("je %s", end_label))

		cg.pushScope()

		for _, stmt := range n.Body {
			cg.GenStmt(stmt)
		}

		cg.popScope()

		cg.loopStartStack = cg.loopStartStack[:len(cg.loopStartStack)-1]
		cg.loopEndStack = cg.loopEndStack[:len(cg.loopEndStack)-1]

		cg.EmitIndent(1, fmt.Sprintf("jmp %s", start_label))

		cg.Emit(fmt.Sprintf("%s:", end_label))

	case *parser.AssignmentStmt:
		val := cg.GenExpr(n.Value)

		offset, ok := cg.lookupVar(n.Name.Name)
		if !ok {
			panic(fmt.Sprintf("undefined variable: %s", n.Name.Name))
		}

		cg.EmitIndent(1, fmt.Sprintf("mov [rbp-%d], %s", offset, val))

	case *parser.BreakStmt:
		if len(cg.loopEndStack) == 0 {
			panic("break statement not inside loop")
		}
		cg.EmitIndent(1, fmt.Sprintf("jmp %s", cg.loopEndStack[len(cg.loopEndStack)-1]))

	case *parser.ContinueStmt:
		if len(cg.loopStartStack) == 0 {
			panic("continue statement not inside loop")
		}
		cg.EmitIndent(1, fmt.Sprintf("jmp %s", cg.loopStartStack[len(cg.loopStartStack)-1]))

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
		offset, ok := cg.lookupVar(n.Name)
		if !ok {
			panic(fmt.Sprintf("undefined variable: %s", n.Name))
		}

		cg.EmitIndent(1, fmt.Sprintf("mov %s, [rbp-%d]", target, offset))
		return target
	case *parser.BoolLit:
		val := 0
		if n.Value {
			val = 1
		}
		cg.EmitIndent(1, fmt.Sprintf("mov %s, %d", target, val))
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

		cg.genExprWithTarget(n.Left, "rax")

		// Allocate temporary stack slot for LHS
		cg.EmitIndent(1, "sub rsp, 8")
		cg.stackPos += 8
		lhsOffset := cg.stackPos
		cg.EmitIndent(1, fmt.Sprintf("mov QWORD [rbp-%d], rax", lhsOffset))

		// Generate right-hand side into rax
		cg.genExprWithTarget(n.Right, "rax")

		// Load LHS back into rbx
		cg.EmitIndent(1, fmt.Sprintf("mov rbx, [rbp-%d]", lhsOffset))

		switch op {
		case "+":
			cg.EmitIndent(1, "add rax, rbx")
		case "-":
			cg.EmitIndent(1, "sub rbx, rax")
			cg.EmitIndent(1, "mov rax, rbx")
		case "*":
			cg.EmitIndent(1, "imul rax, rbx")
		case "<":
			cg.EmitIndent(1, "cmp rbx, rax")
			cg.EmitIndent(1, "setl al")
			cg.EmitIndent(1, "movzx rax, al")
		case ">":
			cg.EmitIndent(1, "cmp rbx, rax")
			cg.EmitIndent(1, "setg al")
			cg.EmitIndent(1, "movzx rax, al")
		case "==":
			cg.EmitIndent(1, "cmp rbx, rax")
			cg.EmitIndent(1, "sete al")
			cg.EmitIndent(1, "movzx rax, al")
		case "<=":
			cg.EmitIndent(1, "cmp rbx, rax")
			cg.EmitIndent(1, "setle al")
			cg.EmitIndent(1, "movzx rax, al")
		case ">=":
			cg.EmitIndent(1, "cmp rbx, rax")
			cg.EmitIndent(1, "setge al")
			cg.EmitIndent(1, "movzx rax, al")
		}

		cg.EmitIndent(1, "add rsp, 8")
		cg.stackPos -= 8

		return target
	default:
		panic(fmt.Sprintf("unsupported expression: %T", n))
	}

}
