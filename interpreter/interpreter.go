package interpreter

import (
	"fmt"
	"strconv"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (intr *Interpreter) Interpret(expr ast.Expr) (any, error) {
	return intr.evaluate(expr)
}

func (intr *Interpreter) PrintEvaluation(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64) // no .0 needed at end
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "nil"
	default:
		e := fmt.Sprintf("Unknown value type evaluated by interpreter: %v", v)
		panic(e)
	}
}

func (intr *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(intr)
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

func (intr *Interpreter) VisitLiteral(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (intr *Interpreter) VisitGrouping(expr *ast.Grouping) (any, error) {
	return intr.evaluate(expr)
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
		rNum, isNumR := right.(float64)
		lNum, isNumL := left.(float64)
		if isNumR && isNumL {
			return rNum + lNum, nil
		}
		// string concatenation
		rStr, isStrR := right.(string)
		lStr, isStrL := left.(string)
		if isStrR && isStrL {
			return rStr + lStr, nil
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
	case lexer.EQUAL:
		return intr.isEqual(left, right), nil
	case lexer.BANG_EQUAL:
		return !intr.isEqual(left, right), nil
	}

	// should be unreachable
	return nil, nil
}
