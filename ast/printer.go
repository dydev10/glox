package ast

import (
	"strings"

	"github.com/dydev10/glox/lexer"
)

type Printer struct {
}

func (p Printer) VisitBinary(expr *Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p Printer) VisitGrouping(expr *Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p Printer) VisitLiteral(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return lexer.PrintLiteral(expr.Value), nil
}

func (p Printer) VisitUnary(expr *Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p Printer) Print(expr Expr) string {
	val, _ := expr.Accept(p)
	return val.(string)
}

func (p Printer) parenthesize(name string, exprs ...Expr) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		val, _ := expr.Accept(p)
		builder.WriteString(val.(string))
	}
	builder.WriteString(")")

	return builder.String(), nil
}
