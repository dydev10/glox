package interpreter

import (
	"fmt"
	"strconv"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[ast.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)

	globals.define("clock", &Clock{})

	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[ast.Expr]int),
	}
}

func PrintEvaluation(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64) // no .0 needed at end
	case bool:
		return strconv.FormatBool(v)
	case fmt.Stringer:
		return v.String()
	case nil:
		return "nil"
	default:
		e := fmt.Sprintf("Unknown value type evaluated by interpreter: %v", v)
		panic(e)
	}
}

// main entry point to run glox statements
func (intr *Interpreter) Interpret(statements []ast.Stmt) error {
	for _, stmt := range statements {
		_, err := intr.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

// alternate entry point to evaluate single expression without statement
func (intr *Interpreter) EvaluateExpression(expr ast.Expr) (any, error) {
	return intr.evaluate(expr)
}

func (intr *Interpreter) execute(stmt ast.Stmt) (any, error) {
	return stmt.Accept(intr)
}

func (intr *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(intr)
}

func (intr *Interpreter) resolve(expr ast.Expr, depth int) {
	intr.locals[expr] = depth
}

func (intr *Interpreter) executeBlock(statements []ast.Stmt, env *Environment) error {
	prevEnv := intr.environment
	intr.environment = env

	for _, stmt := range statements {
		_, err := intr.execute(stmt)
		if err != nil {
			// restore original environment on error
			intr.environment = prevEnv
			return err
		}
	}

	// restore original environment before returning
	intr.environment = prevEnv
	return nil
}

func (intr *Interpreter) lookupVariable(name *lexer.Token, expr ast.Expr) (any, error) {
	distance, ok := intr.locals[expr]
	if ok {
		return intr.environment.getAt(distance, name.Lexeme)
	} else {
		return intr.globals.get(name)
	}
}

func (intr *Interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if v, ok := val.(bool); ok {
		return v
	}
	return true
}

func (intr *Interpreter) isEqual(a, b any) bool {
	// first check if both nil, both nil means equal
	if a == nil && b == nil {
		return true
	}
	// check if both type and value match
	// number match
	aNum, isNumA := a.(float64)
	bNum, isNumB := b.(float64)
	if isNumA && isNumB {
		return aNum == bNum
	}
	// string equality
	aStr, isStrA := a.(string)
	bStr, isStrB := b.(string)
	if isStrA && isStrB {
		return aStr == bStr
	}
	// bool equality
	aBool, isBoolA := a.(bool)
	bBool, isBoolB := b.(bool)
	if isBoolA && isBoolB {
		return aBool == bBool
	}
	// type mismatch, mean inequality
	return false
}

func checkNumberOperand(operator *lexer.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}

	return &RuntimeError{token: operator, message: "Operand must be a number."}
}

func checkNumberOperands(operator *lexer.Token, left any, right any) error {
	_, okL := left.(float64)
	_, okR := right.(float64)
	if okL && okR {
		return nil
	}

	return &RuntimeError{token: operator, message: "Operands must be numbers."}
}

/*
* Expr interface implementation
 */

func (intr *Interpreter) VisitLiteral(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (intr *Interpreter) VisitLogical(expr *ast.Logical) (any, error) {
	left, err := intr.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == lexer.OR {
		if intr.isTruthy(left) {
			return left, nil
		}
	}
	if expr.Operator.Type == lexer.AND {
		if !intr.isTruthy(left) {
			return left, nil
		}
	}

	return intr.evaluate(expr.Right)
}

func (intr *Interpreter) VisitSet(expr *ast.Set) (any, error) {
	object, objErr := intr.evaluate(expr.Object)
	if objErr != nil {
		return nil, objErr
	}

	instance, isInstance := object.(*LoxInstance)
	if !isInstance {
		return nil, &RuntimeError{
			token:   expr.Name,
			message: "Only instances have fields.",
		}
	}

	value, valueErr := intr.evaluate(expr.Value)
	if valueErr != nil {
		return nil, valueErr
	}

	instance.Set(expr.Name, value)

	return value, nil
}

func (intr *Interpreter) VisitThis(expr *ast.This) (any, error) {
	return intr.lookupVariable(expr.Keyword, expr)
}

func (intr *Interpreter) VisitGrouping(expr *ast.Grouping) (any, error) {
	return intr.evaluate(expr.Expression)
}

func (intr *Interpreter) VisitUnary(expr *ast.Unary) (any, error) {
	right, err := intr.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case lexer.MINUS:
		err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case lexer.BANG:
		return !intr.isTruthy(right), nil
	}

	// should be unreachable
	return nil, nil
}

func (intr *Interpreter) VisitVariable(expr *ast.Variable) (any, error) {
	return intr.lookupVariable(expr.Name, expr)
}

func (intr *Interpreter) VisitBinary(expr *ast.Binary) (any, error) {
	left, lErr := intr.evaluate(expr.Left)
	if lErr != nil {
		return nil, lErr
	}
	right, rErr := intr.evaluate(expr.Right)
	if rErr != nil {
		return nil, rErr
	}

	switch expr.Operator.Type {
	// arithmetic
	case lexer.MINUS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case lexer.SLASH:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case lexer.STAR:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	// overloaded arithmetic / string concatenation
	case lexer.PLUS:
		// number addition
		lNum, isNumL := left.(float64)
		rNum, isNumR := right.(float64)
		if isNumL && isNumR {
			return lNum + rNum, nil
		}
		// string concatenation
		lStr, isStrL := left.(string)
		rStr, isStrR := right.(string)
		if isStrL && isStrR {
			return lStr + rStr, nil
		}
		// type match failed, return error
		return nil, &RuntimeError{
			token:   expr.Operator,
			message: "Operands must be two numbers or two strings.",
		}

	// comparison
	case lexer.GREATER:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case lexer.GREATER_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case lexer.LESS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case lexer.LESS_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil

	// equality
	case lexer.EQUAL_EQUAL:
		return intr.isEqual(left, right), nil
	case lexer.BANG_EQUAL:
		return !intr.isEqual(left, right), nil
	}

	// should be unreachable
	return nil, nil
}

func (intr *Interpreter) VisitCall(expr *ast.Call) (any, error) {
	callee, calleeErr := intr.evaluate(expr.Callee)
	if calleeErr != nil {
		return nil, calleeErr
	}

	arguments := []any{}
	for _, arg := range expr.Arguments {
		argEval, argEvalErr := intr.evaluate(arg)
		if argEvalErr != nil {
			return nil, argEvalErr
		}

		arguments = append(arguments, argEval)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		notCallableErr := &RuntimeError{
			token:   expr.Paren,
			message: "Can only call functions and classes.",
		}
		return nil, notCallableErr
	}

	if function.Arity() != len(arguments) {
		arityErr := &RuntimeError{
			token:   expr.Paren,
			message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)),
		}
		return nil, arityErr
	}

	return function.Call(intr, arguments)
}

func (intr *Interpreter) VisitGet(expr *ast.Get) (any, error) {
	object, objectErr := intr.evaluate(expr.Object)
	if objectErr != nil {
		return nil, objectErr
	}

	instance, isInstance := object.(*LoxInstance)
	if !isInstance {
		return &RuntimeError{
			token:   expr.Name,
			message: "Only instances have properties.",
		}, nil
	}

	return instance.Get(expr.Name)
}

func (intr *Interpreter) VisitAssign(expr *ast.Assign) (any, error) {
	value, err := intr.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := intr.locals[expr]
	if ok {
		intr.environment.assignAt(distance, expr.Name, value)
	} else {
		assignErr := intr.globals.assign(expr.Name, value)
		if assignErr != nil {
			return nil, assignErr
		}
	}

	return value, nil
}

/*
* Stmt interface implementation
 */

func (intr *Interpreter) VisitBlock(stmt *ast.Block) (any, error) {
	err := intr.executeBlock(stmt.Statements, NewEnvironment(intr.environment))
	return nil, err
}

func (intr *Interpreter) VisitClass(stmt *ast.Class) (any, error) {
	var superclass *LoxClass
	if stmt.Superclass != nil {
		class, err := intr.evaluate(stmt.Superclass)
		if err != nil {
			return nil, err
		}
		loxSuperclass, isLoxClass := class.(*LoxClass)
		if !isLoxClass {
			return nil, &RuntimeError{
				token:   stmt.Superclass.Name,
				message: "Superclass must be a class.",
			}
		}
		superclass = loxSuperclass
	}

	intr.environment.define(stmt.Name.Lexeme, nil)

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		function := &LoxFunction{
			declaration:   method,
			closure:       intr.environment,
			isInitializer: method.Name.Lexeme == "init", // check if method is constructor
		}
		methods[method.Name.Lexeme] = function
	}

	class := &LoxClass{
		name:       stmt.Name.Lexeme,
		superclass: superclass,
		methods:    methods,
	}
	intr.environment.assign(stmt.Name, class)

	return nil, nil
}

func (intr *Interpreter) VisitExpression(stmt *ast.Expression) (any, error) {
	_, err := intr.evaluate(stmt.Expression)
	return nil, err
}

func (intr *Interpreter) VisitFunction(stmt *ast.Function) (any, error) {
	function := &LoxFunction{
		declaration:   stmt,
		closure:       intr.environment,
		isInitializer: false,
	}
	intr.environment.define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (intr *Interpreter) VisitIf(stmt *ast.If) (any, error) {
	condResult, condErr := intr.evaluate(stmt.Condition)
	if condErr != nil {
		return nil, condErr
	}

	var execErr error
	if intr.isTruthy(condResult) {
		_, execErr = intr.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		_, execErr = intr.execute(stmt.ElseBranch)
	}

	return nil, execErr
}

func (intr *Interpreter) VisitPrint(stmt *ast.Print) (any, error) {
	val, err := intr.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", PrintEvaluation(val))
	return nil, nil
}

func (intr *Interpreter) VisitReturn(stmt *ast.Return) (any, error) {
	var value any
	var valErr error

	if stmt.Value != nil {
		value, valErr = intr.evaluate(stmt.Value)
		if valErr != nil {
			return nil, valErr
		}
	}

	// intentionally return value as error so that it can skip the call stack to return control to LoxFunction.Call
	return nil, &ThrownReturn{value: value}
}

func (intr *Interpreter) VisitVar(stmt *ast.Var) (any, error) {
	var value any
	var err error

	if stmt.Initializer != nil {
		value, err = intr.evaluate(stmt.Initializer)
	}

	intr.environment.define(stmt.Name.Lexeme, value)
	return nil, err
}

func (intr *Interpreter) VisitWhile(stmt *ast.While) (any, error) {
	cond, condErr := intr.evaluate(stmt.Condition)
	if condErr != nil {
		return nil, condErr
	}

	for intr.isTruthy(cond) {
		_, err := intr.execute(stmt.Body)
		if err != nil {
			return nil, err
		}

		cond, condErr = intr.evaluate(stmt.Condition)
		if condErr != nil {
			return nil, condErr
		}
	}

	return nil, nil
}
