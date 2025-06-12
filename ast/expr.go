package ast

import "github.com/dydev10/glox/lexer"

type Expr interface {
	Accept(v VisitorExpr[any]) (any, error)
}

type VisitorExpr[R any] interface {
	VisitAssign(expr *Assign) (R, error)
	VisitBinary(expr *Binary) (R, error)
	VisitGrouping(expr *Grouping) (R, error)
	VisitLiteral(expr *Literal) (R, error)
	VisitUnary(expr *Unary) (R, error)
	VisitVariable(expr *Variable) (R, error)
}

type Assign struct {
	Name  *lexer.Token
	Value Expr
}

func (n *Assign) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitAssign(n)
}

type Binary struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func (n *Binary) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitBinary(n)
}

type Grouping struct {
	Expression Expr
}

func (n *Grouping) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitGrouping(n)
}

type Literal struct {
	Value any
}

func (n *Literal) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitLiteral(n)
}

type Unary struct {
	Operator *lexer.Token
	Right    Expr
}

func (n *Unary) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitUnary(n)
}

type Variable struct {
	Name *lexer.Token
}

func (n *Variable) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitVariable(n)
}
