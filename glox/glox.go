package glox

import (
	"fmt"
	"os"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/interpreter"
	"github.com/dydev10/glox/lexer"
	"github.com/dydev10/glox/parser"
)

type Glox struct {
	source     string
	command    string
	isRunMode  bool
	isEvalMode bool

	HadSyntaxError  bool
	HadRuntimeError bool
	errorList       []error

	tokens     []*lexer.Token
	expression ast.Expr
	statements []ast.Stmt

	evaluation any
}

func NewGlox(command, source string) *Glox {
	return &Glox{
		source:     source,
		command:    command,
		isRunMode:  command == "run",
		isEvalMode: command == "evaluate",
	}
}

func (g *Glox) Tokenize() {
	//lexer
	l := lexer.New(g.source)
	g.tokens = l.Lex()
	if len(l.Errors) > 0 {
		g.HadSyntaxError = true
		g.errorList = append(g.errorList, l.Errors...)
	}
}

func (g *Glox) RunStatements() {
	if g.HadSyntaxError {
		return
	}

	p := parser.NewParser(g.tokens)
	statements, parseError := p.Parse()
	g.statements = statements
	if parseError != nil {
		g.HadSyntaxError = true
		g.errorList = append(g.errorList, parseError)
	}

	// end execution if only parse command
	if !g.isRunMode || g.HadSyntaxError {
		return
	}
	intr := interpreter.NewInterpreter()
	runtimeErr := intr.Interpret(statements)
	if runtimeErr != nil {
		g.HadRuntimeError = true
		g.errorList = append(g.errorList, runtimeErr)
	}
}

func (g *Glox) RunExpression() {
	if g.HadSyntaxError {
		return
	}

	p := parser.NewParser(g.tokens)
	expression, parseError := p.ParseExpression()
	g.expression = expression
	if parseError != nil {
		g.HadSyntaxError = true
		g.errorList = append(g.errorList, parseError)
	}

	// end execution if only parse command
	if !g.isEvalMode || g.HadSyntaxError {
		return
	}
	intr := interpreter.NewInterpreter()
	evaluation, runtimeErr := intr.EvaluateExpression(expression)
	g.evaluation = evaluation
	if runtimeErr != nil {
		g.HadRuntimeError = true
		g.errorList = append(g.errorList, runtimeErr)
	}
}

func (g *Glox) PrintErrors() {
	for _, lexError := range g.errorList {
		fmt.Fprintln(os.Stderr, lexError)
	}
}

func (g *Glox) PrintResult() {
	switch g.command {
	case "tokenize":
		for _, v := range g.tokens {
			fmt.Println(v.String())
		}
	case "parse":
		if !g.HadSyntaxError {
			astPrinter := &ast.Printer{}
			out := astPrinter.Print(g.expression)
			fmt.Printf("%s", out)
		}
	case "evaluate":
		if !g.HadSyntaxError && !g.HadRuntimeError {
			evalOut := interpreter.PrintEvaluation(g.evaluation)
			fmt.Printf("%s", evalOut)
		}
	}
}
