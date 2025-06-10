package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dydev10/glox/glox"
)

func main() {
	if len(os.Args) == 1 {
		startREPL()
		os.Exit(0)
	} else if len(os.Args) == 3 {
		runFile(os.Args[1], os.Args[2])
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}
}

func startREPL() {
	print("Welcome to glox!\nEnter expression to evaluate.\n")
	defer os.Exit(0)

	REPL()
}

func runFile(command, filename string) {
	if command != "tokenize" && command != "parse" && command != "evaluate" && command != "run" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	glox := glox.NewGlox(command, string(fileContents))
	glox.Tokenize()

	if command == "parse" || command == "evaluate" {
		glox.RunExpression()
	}

	if command == "run" {
		glox.RunStatements()
	}

	glox.PrintErrors()
	glox.PrintResult()

	if glox.HadRuntimeError {
		os.Exit(70)
	} else if glox.HadSyntaxError {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}

func REPL() {
	print("\n> ")

	// read input line
	bufScanner := bufio.NewScanner(os.Stdin)
	bufScanner.Scan()
	source := strings.TrimSpace(bufScanner.Text())
	bufScanner = nil
	if source == "exit" {
		return
	} else {
		defer REPL() // run in loop on returning unless exit input
	}

	var command string
	if strings.HasSuffix(source, ";") {
		command = "run"
	} else {
		command = "evaluate"
	}

	glox := glox.NewGlox(command, source)
	glox.Tokenize()

	if command == "run" {
		glox.RunStatements()
	} else {
		glox.RunExpression()
	}

	glox.PrintErrors()
	glox.PrintResult()
}
