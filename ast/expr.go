package ast

import "github.com/dydev10/glox/lexer"

type Expr interface {
	Accept(v Visitor[any]) any
}

type Visitor[R any] interface {
	VisitBinary(expr *Binary) R
	VisitGrouping(expr *Grouping) R
	VisitLiteral(expr *Literal) R
	VisitUnary(expr *Unary) R
}

type Binary struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func (n *Binary) Accept(v Visitor[any]) any {
	return v.VisitBinary(n)
}

type Grouping struct {
	Expression Expr
}

func (n *Grouping) Accept(v Visitor[any]) any {
	return v.VisitGrouping(n)
}

type Literal struct {
	Value any
}

func (n *Literal) Accept(v Visitor[any]) any {
	return v.VisitLiteral(n)
}

type Unary struct {
	Operator *lexer.Token
	Right    Expr
}

func (n *Unary) Accept(v Visitor[any]) any {
	return v.VisitUnary(n)
}
