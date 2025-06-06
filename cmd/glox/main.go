package main

import (
	"fmt"
	"os"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
	"github.com/dydev10/glox/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	/*
	* Main Lexer Code
	 */
	l := lexer.New(string(fileContents))
	tokens := l.Lex()
	hadErrors := len(l.Errors) > 0

	if hadErrors {
		hadErrors = true
		for _, lexError := range l.Errors {
			fmt.Fprintln(os.Stderr, lexError.String())
		}
	}

	if command == "tokenize" {
		for _, v := range tokens {
			fmt.Println(v.String())
		}
	}

	if command == "parse" {
		p := parser.NewParser(tokens)
		expression := p.Parse()
		printer := &ast.Printer{}
		out := printer.Print(expression)
		fmt.Printf("%s", out)
	}

	if hadErrors {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}
