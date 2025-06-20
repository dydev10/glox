package ast

import (
	"strings"

	"github.com/dydev10/glox/lexer"
)

type Printer struct {
}

func (p Printer) VisitAssign(expr *Assign) (any, error) {
	// TODO: check if this works fine?? probably not
	return p.parenthesize("=", expr, expr.Value)
}

func (p Printer) VisitBinary(expr *Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p Printer) VisitCall(expr *Call) (any, error) {
	var builder strings.Builder

	callee, _ := expr.Callee.Accept(p)
	builder.WriteString(callee.(string))
	builder.WriteString("(")
	for i, arg := range expr.Arguments {
		val, _ := arg.Accept(p)
		builder.WriteString(val.(string))
		if i < len(expr.Arguments)-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString(")")

	return builder.String(), nil
}

func (p Printer) VisitGet(expr *Get) (any, error) {
	// TODO: improve print representation of this tree node
	return p.parenthesize("get "+expr.Name.Lexeme, expr.Object)
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

func (p Printer) VisitLogical(expr *Logical) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p Printer) VisitSet(expr *Set) (any, error) {
	// TODO: improve print representation of this tree node
	return p.parenthesize("set "+expr.Name.Lexeme, expr.Object, expr.Value)
}

func (p Printer) VisitUnary(expr *Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p Printer) VisitVariable(expr *Variable) (any, error) {
	return expr.Name.Lexeme, nil
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
