package ast

import "github.com/dydev10/glox/lexer"

type Expr[R any] interface {
	Accept(v Visitor[R]) R
}

type Visitor[R any] interface {
	VisitBinary(expr *Binary[R]) R
	VisitGrouping(expr *Grouping[R]) R
	VisitLiteral(expr *Literal[R]) R
	VisitUnary(expr *Unary[R]) R
}


type Binary[R any] struct {
	Left Expr[R]
	Operator *lexer.Token
	Right Expr[R]
}

func (n *Binary[R]) Accept(v Visitor[R]) R {
	return v.VisitBinary(n)
}

type Grouping[R any] struct {
	Expression Expr[R]
}

func (n *Grouping[R]) Accept(v Visitor[R]) R {
	return v.VisitGrouping(n)
}

type Literal[R any] struct {
	Value any
}

func (n *Literal[R]) Accept(v Visitor[R]) R {
	return v.VisitLiteral(n)
}

type Unary[R any] struct {
	Operator *lexer.Token
	Right Expr[R]
}

func (n *Unary[R]) Accept(v Visitor[R]) R {
	return v.VisitUnary(n)
}

