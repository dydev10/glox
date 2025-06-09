package ast

type Stmt interface {
	Accept(v VisitorStmt[any]) (any, error)
}

type VisitorStmt[R any] interface {
	VisitExpression(expr *Expression) (R, error)
	VisitPrint(expr *Print) (R, error)
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
