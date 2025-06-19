package ast

import "github.com/dydev10/glox/lexer"

type Stmt interface {
	Accept(v VisitorStmt[any]) (any, error)
}

type VisitorStmt[R any] interface {
	VisitBlock(expr *Block) (R, error)
	VisitClass(expr *Class) (R, error)
	VisitExpression(expr *Expression) (R, error)
	VisitFunction(expr *Function) (R, error)
	VisitIf(expr *If) (R, error)
	VisitPrint(expr *Print) (R, error)
	VisitReturn(expr *Return) (R, error)
	VisitVar(expr *Var) (R, error)
	VisitWhile(expr *While) (R, error)
}

type Block struct {
	Statements []Stmt
}

func (n *Block) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitBlock(n)
}

type Class struct {
	Name    *lexer.Token
	Methods []*Function
}

func (n *Class) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitClass(n)
}

type Expression struct {
	Expression Expr
}

func (n *Expression) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitExpression(n)
}

type Function struct {
	Name   *lexer.Token
	Params []*lexer.Token
	Body   []Stmt
}

func (n *Function) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitFunction(n)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (n *If) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitIf(n)
}

type Print struct {
	Expression Expr
}

func (n *Print) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitPrint(n)
}

type Return struct {
	Keyword *lexer.Token
	Value   Expr
}

func (n *Return) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitReturn(n)
}

type Var struct {
	Name        *lexer.Token
	Initializer Expr
}

func (n *Var) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitVar(n)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (n *While) Accept(v VisitorStmt[any]) (any, error) {
	return v.VisitWhile(n)
}
