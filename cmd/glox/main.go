package main

import (
	"fmt"
	"os"

	"github.com/dydev10/glox/lexer"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
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

	for _, v := range tokens {
		fmt.Println(v.String())
	}

	if hadErrors {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}
