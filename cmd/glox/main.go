package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/interpreter"
	"github.com/dydev10/glox/lexer"
	"github.com/dydev10/glox/parser"
)

func main() {
	if len(os.Args) == 1 {
		run()
		os.Exit(0)
	} else if len(os.Args) == 3 {
		runFile(os.Args[1], os.Args[2])
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}
}

func run() {
	print("Welcome to glox!\nEnter expression to evaluate.\n")
	defer os.Exit(0)

	REPL()
}

func runFile(command, filename string) {
	if command != "tokenize" && command != "parse" && command != "evaluate" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

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

	// interpreter
	var eval any
	var runtimeErr error
	hadRuntimeErrors := false
	intr := interpreter.NewInterpreter()
	if !hadParseErrors {
		eval, runtimeErr = intr.Interpret(expression)
	}
	if runtimeErr != nil {
		hadRuntimeErrors = true
	}

	// printing error
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

	if command == "evaluate" && hadRuntimeErrors {
		hadErrors = true
		fmt.Fprintln(os.Stderr, runtimeErr)
	}

	// printing output
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

	if command == "evaluate" && !hadRuntimeErrors {
		evalOut := intr.PrintEvaluation(eval)
		fmt.Printf("%s", evalOut)
	}

	if hadRuntimeErrors {
		os.Exit(70)
	} else if hadErrors {
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
	source := bufScanner.Text()
	bufScanner = nil
	if source == "exit" {
		return
	} else {
		defer REPL() // run in loop on returning unless exit input
	}

	//lexer
	l := lexer.New(source)
	tokens := l.Lex()
	hadLexErrors := len(l.Errors) > 0
	if hadLexErrors {
		for _, lexError := range l.Errors {
			fmt.Fprintln(os.Stderr, lexError.String())
		}
		return
	}

	// parser
	p := parser.NewParser(tokens)
	expression, parseError := p.Parse()
	if parseError != nil {
		for _, parseError := range p.Errors {
			fmt.Fprintln(os.Stderr, parseError.String())
		}
		return
	}

	// interpreter
	intr := interpreter.NewInterpreter()
	eval, runtimeErr := intr.Interpret(expression)
	if runtimeErr != nil {
		fmt.Fprintln(os.Stderr, runtimeErr)
		return
	}
	evalOut := intr.PrintEvaluation(eval)
	fmt.Println(evalOut)
}
