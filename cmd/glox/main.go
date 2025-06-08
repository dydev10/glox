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
	hadErrors := false
	//lexer
	l := lexer.New(string(fileContents))
	tokens := l.Lex()
	hadLexErrors := len(l.Errors) > 0

	// parser
	p := parser.NewParser(tokens)
	expression, parseError := p.Parse()
	hadParseErrors := false
	if parseError != nil {
		hadParseErrors = true
	}

	if hadLexErrors {
		hadErrors = true
		for _, lexError := range l.Errors {
			fmt.Fprintln(os.Stderr, lexError.String())
		}
	}

	if command == "parse" && hadParseErrors {
		hadErrors = true
		for _, parseError := range p.Errors {
			fmt.Fprintln(os.Stderr, parseError.String())
		}
	}

	if command == "tokenize" {
		for _, v := range tokens {
			fmt.Println(v.String())
		}
	}

	if command == "parse" && !hadParseErrors {
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
