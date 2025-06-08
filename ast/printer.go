package ast

import (
	"strings"

	"github.com/dydev10/glox/lexer"
)

type Printer struct {
}

func (p Printer) VisitBinary(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p Printer) VisitGrouping(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p Printer) VisitLiteral(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return lexer.PrintLiteral(expr.Value)
}

func (p Printer) VisitUnary(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p Printer) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p Printer) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p).(string))
	}
	builder.WriteString(")")

	return builder.String()
}
