package ast

import "github.com/dydev10/glox/lexer"

type Expr interface {
	Accept(v Visitor[any]) (any, error)
}

type Visitor[R any] interface {
	VisitBinary(expr *Binary) (R, error)
	VisitGrouping(expr *Grouping) (R, error)
	VisitLiteral(expr *Literal) (R, error)
	VisitUnary(expr *Unary) (R, error)
}

type Binary struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func (n *Binary) Accept(v Visitor[any]) (any, error) {
	return v.VisitBinary(n)
}

type Grouping struct {
	Expression Expr
}

func (n *Grouping) Accept(v Visitor[any]) (any, error) {
	return v.VisitGrouping(n)
}

type Literal struct {
	Value any
}

func (n *Literal) Accept(v Visitor[any]) (any, error) {
	return v.VisitLiteral(n)
}

type Unary struct {
	Operator *lexer.Token
	Right    Expr
}

func (n *Unary) Accept(v Visitor[any]) (any, error) {
	return v.VisitUnary(n)
}
