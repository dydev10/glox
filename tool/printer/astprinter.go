package main

import (
	"fmt"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
)

func main() {
	expression := &ast.Binary{
		Left: &ast.Unary{
			Operator: &lexer.Token{
				Type:    lexer.MINUS,
				Lexeme:  "-",
				Literal: "",
			},
			Right: &ast.Literal{
				Value: float64(123),
			},
		},
		Operator: &lexer.Token{
			Type:    lexer.STAR,
			Lexeme:  "*",
			Literal: "",
		},
		Right: &ast.Grouping{
			Expression: &ast.Literal{
				Value: float64(45.67),
			},
		},
	}

	printer := &ast.Printer{}
	out := printer.Print(expression)

	fmt.Printf("%s", out)
}
