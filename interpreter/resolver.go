package interpreter

import (
	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/ds"
	"github.com/dydev10/glox/lexer"
)

type BlockScope map[string]bool

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          *ds.Stack[BlockScope]
	currentFunction FunctionType
	Errors          []error
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          ds.NewStack[BlockScope](),
		currentFunction: NONE,
		Errors:          []error{},
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) {
	r.resolveStatements(statements)
}

func (r *Resolver) logError(token *lexer.Token, message string) {
	r.Errors = append(r.Errors, &ResolveError{
		token:   token,
		message: message,
	})
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveStatements(statements []ast.Stmt) {
	for _, stmt := range statements {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(BlockScope))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name *lexer.Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek()
	if _, alreadyDeclared := scope[name.Lexeme]; alreadyDeclared {
		r.logError(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *lexer.Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name *lexer.Token) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		if _, ok := r.scopes.Get(i)[name.Lexeme]; ok {
			// TODO: handle the possible error returned by resolve, if it has error as return value
			r.interpreter.resolve(expr, r.scopes.Len()-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *ast.Function, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

/**
*	Stmt Visitor interface implementation
 */

func (r *Resolver) VisitBlock(stmt *ast.Block) (any, error) {
	r.beginScope()
	r.resolveStatements(stmt.Statements)
	r.endScope()

	return nil, nil
}

func (r *Resolver) VisitClass(stmt *ast.Class) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	return nil, nil
}

func (r *Resolver) VisitVar(stmt *ast.Var) (any, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)

	return nil, nil
}

func (r *Resolver) VisitFunction(stmt *ast.Function) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUNCTION)

	return nil, nil
}

func (r *Resolver) VisitExpression(stmt *ast.Expression) (any, error) {
	r.resolveExpr(stmt.Expression)

	return nil, nil
}

func (r *Resolver) VisitIf(stmt *ast.If) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}

	return nil, nil
}

func (r *Resolver) VisitPrint(stmt *ast.Print) (any, error) {
	r.resolveExpr(stmt.Expression)

	return nil, nil
}

func (r *Resolver) VisitReturn(stmt *ast.Return) (any, error) {
	if r.currentFunction == NONE {
		r.logError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil, nil
}

func (r *Resolver) VisitWhile(stmt *ast.While) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)

	return nil, nil
}

/**
*	Expr Visitor interface implementation
 */

func (r *Resolver) VisitVariable(expr *ast.Variable) (any, error) {
	if !r.scopes.IsEmpty() {
		defined, declared := r.scopes.Peek()[expr.Name.Lexeme]
		if declared && !defined {
			r.logError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitAssign(expr *ast.Assign) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitBinary(expr *ast.Binary) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil, nil
}

func (r *Resolver) VisitCall(expr *ast.Call) (any, error) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGrouping(expr *ast.Grouping) (any, error) {
	r.resolveExpr(expr.Expression)

	return nil, nil
}

func (r *Resolver) VisitLiteral(expr *ast.Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogical(expr *ast.Logical) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil, nil
}

func (r *Resolver) VisitUnary(expr *ast.Unary) (any, error) {
	r.resolveExpr(expr.Right)

	return nil, nil
}
