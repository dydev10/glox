package interpreter

import (
	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/ds"
	"github.com/dydev10/glox/lexer"
)

type BlockScope map[string]bool

type FunctionType int

const (
	ftNONE FunctionType = iota
	ftFUNCTION
	ftINITIALIZER
	ftMETHOD
)

type ClassType int

const (
	ctNONE ClassType = iota
	ctCLASS
	ctSUBCLASS
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          *ds.Stack[BlockScope]
	currentFunction FunctionType
	currentClass    ClassType
	Errors          []error
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          ds.NewStack[BlockScope](),
		currentFunction: ftNONE,
		currentClass:    ctNONE,
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
	enclosingClass := r.currentClass
	r.currentClass = ctCLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		r.logError(stmt.Superclass.Name, "A class can't inherit from itself.")
	}

	if stmt.Superclass != nil {
		r.currentClass = ctSUBCLASS
		r.resolveExpr(stmt.Superclass)

		// inject new scope to resolve super keyword for this class's methods
		r.beginScope()
		r.scopes.Peek()["super"] = true
	}

	r.beginScope()
	r.scopes.Peek()["this"] = true

	for _, method := range stmt.Methods {
		declaration := ftMETHOD
		if method.Name.Lexeme == "init" {
			declaration = ftINITIALIZER
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()
	if stmt.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass

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
	r.resolveFunction(stmt, ftFUNCTION)

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
	if r.currentFunction == ftNONE {
		r.logError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		// only block returning value from constructor. allow empty return for early exits
		if r.currentFunction == ftINITIALIZER {
			r.logError(stmt.Keyword, "Can't return a value from an initializer.")
		}

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

func (r *Resolver) VisitGet(expr *ast.Get) (any, error) {
	r.resolveExpr(expr.Object)

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

func (r *Resolver) VisitThis(expr *ast.This) (any, error) {
	if r.currentClass == ctNONE {
		r.logError(expr.Keyword, "Can't use 'this' outside of a class.")
	}

	r.resolveLocal(expr, expr.Keyword)

	return nil, nil
}

func (r *Resolver) VisitSet(expr *ast.Set) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)

	return nil, nil
}

func (r *Resolver) VisitSuper(expr *ast.Super) (any, error) {
	if r.currentClass == ctNONE {
		r.logError(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != ctSUBCLASS {
		r.logError(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	}

	r.resolveLocal(expr, expr.Keyword)

	return nil, nil
}

func (r *Resolver) VisitUnary(expr *ast.Unary) (any, error) {
	r.resolveExpr(expr.Right)

	return nil, nil
}
