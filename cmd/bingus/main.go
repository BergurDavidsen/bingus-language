package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BergurDavidsen/bingus/internal/codegen"
	"github.com/BergurDavidsen/bingus/internal/lexer"
	"github.com/BergurDavidsen/bingus/internal/parser"
)

var file_extension = ".bng"
var output_folder = "./output/"

func generateOutputFiles(asm string) {
	err := os.WriteFile(fmt.Sprintf("%stest.asm", output_folder), []byte(asm), 0644)
	if err != nil {
		panic(fmt.Sprintf("error: %s", err))
	}

	cmd := exec.Command("nasm", "-f", "elf64", fmt.Sprintf("%stest.asm", output_folder), "-o", fmt.Sprintf("%stest.o", output_folder))
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("ld", "-o", fmt.Sprintf("%stest", output_folder), fmt.Sprintf("%stest.o", output_folder))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}
	fmt.Println("Compiled file successfully!")
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: bingus <filename>%s\n", file_extension)
		os.Exit(1)
	}
	filename := os.Args[1]
	if filepath.Ext(filename) != file_extension {
		fmt.Printf("Error: file must have %s extension (got %s)\n", file_extension, filepath.Ext(filename))
		os.Exit(1)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}
	// fmt.Println("Source code:")
	// fmt.Println(string(data))
	// fmt.Println("")

	tokens := lexer.Lex(string(data))
	// fmt.Println("Tokens:")

	// for _, tok := range tokens {
	// 	fmt.Printf("%#v\n", tok)
	// }

	p := parser.Parser{Tokens: tokens}
	program := p.ParseProgram()
	// parser.PrintNodeReflect(program, "")

	// env := NewEnv()
	// result := env.Eval(program)
	// fmt.Println("Result: ", result)

	cg := codegen.NewCodeGen()
	cg.Gen(program)

	asm := cg.String()

	generateOutputFiles(asm)
}
