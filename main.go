package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	file_extension := ".bng"
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
	fmt.Println("Source code:")
	fmt.Println(string(data))
	fmt.Println("")
	tokens := lex(string(data))
	fmt.Println("Tokens:")
	for _, tok := range tokens {
		fmt.Printf("%#v\n", tok)
	}

	fmt.Println("")
	parser := Parser{tokens: tokens}
	program := parser.parseProgram()

	printNodeReflect(program, "")
}
