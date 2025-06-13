package ast

import "github.com/dydev10/glox/lexer"

type Stmt interface {
	Accept(v VisitorStmt[any]) (any, error)
}

type VisitorStmt[R any] interface {
	VisitBlock(expr *Block) (R, error)
	VisitExpression(expr *Expression) (R, error)
	VisitPrint(expr *Print) (R, error)
	VisitVar(expr *Var) (R, error)
}

type Block struct {
	Statements []Stmt
}

func (n *Block) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitBlock(n)
}

type Expression struct {
	Expression Expr
}

func (n *Expression) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitExpression(n)
}

type Print struct {
	Expression Expr
}

func (n *Print) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitPrint(n)
}

type Var struct {
	Name        *lexer.Token
	Initializer Expr
}

func (n *Var) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitVar(n)
}
