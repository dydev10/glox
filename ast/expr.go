package ast

import "github.com/dydev10/glox/lexer"

type Expr interface {
	Accept(v VisitorExpr[any]) (any, error)
}

type VisitorExpr[R any] interface {
	VisitAssign(expr *Assign) (R, error)
	VisitBinary(expr *Binary) (R, error)
	VisitCall(expr *Call) (R, error)
	VisitGet(expr *Get) (R, error)
	VisitGrouping(expr *Grouping) (R, error)
	VisitLiteral(expr *Literal) (R, error)
	VisitLogical(expr *Logical) (R, error)
	VisitSet(expr *Set) (R, error)
	VisitSuper(expr *Super) (R, error)
	VisitThis(expr *This) (R, error)
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

type Call struct {
	Callee    Expr
	Paren     *lexer.Token
	Arguments []Expr
}

func (n *Call) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitCall(n)
}

type Get struct {
	Object Expr
	Name   *lexer.Token
}

func (n *Get) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitGet(n)
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

type Logical struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func (n *Logical) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitLogical(n)
}

type Set struct {
	Object Expr
	Name   *lexer.Token
	Value  Expr
}

func (n *Set) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitSet(n)
}

type Super struct {
	Keyword *lexer.Token
	Method  *lexer.Token
}

func (n *Super) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitSuper(n)
}

type This struct {
	Keyword *lexer.Token
}

func (n *This) Accept(v VisitorExpr[any]) (any, error) {
	return v.VisitThis(n)
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
