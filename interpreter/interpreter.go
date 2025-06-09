package interpreter

import (
	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
)

type Interpreter struct {
}

func (intr *Interpreter) evaluate(expr ast.Expr) any {
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

func (intr *Interpreter) VisitLiteral(expr *ast.Literal) any {
	return expr.Value
}

func (intr *Interpreter) VisitGrouping(expr *ast.Grouping) any {
	return intr.evaluate(expr)
}

func (intr *Interpreter) VisitUnary(expr *ast.Unary) any {
	right := intr.evaluate(expr.Right)
	// handle right error

	switch expr.Operator.Type {
	case lexer.MINUS:
		return -right.(float64)
	case lexer.BANG:
		return !intr.isTruthy(right)
	}

	// should be unreachable
	return nil
}

func (intr *Interpreter) VisitBinary(expr *ast.Binary) any {
	left := intr.evaluate(expr.Left)
	right := intr.evaluate(expr.Right)

	switch expr.Operator.Type {
	// arithmetic
	case lexer.MINUS:
		return left.(float64) - right.(float64)
	case lexer.SLASH:
		return left.(float64) / right.(float64)
	case lexer.STAR:
		return left.(float64) * right.(float64)

		// overloaded arithmetic / string concatenation
	case lexer.PLUS:
		// number addition
		rNum, isNumR := right.(float64)
		lNum, isNumL := left.(float64)
		if isNumR && isNumL {
			return rNum + lNum
		}

		// string concatenation
		rStr, isStrR := right.(string)
		lStr, isStrL := left.(string)
		if isStrR && isStrL {
			return rStr + lStr
		}

	// comparison
	case lexer.GREATER:
		return left.(float64) > right.(float64)
	case lexer.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case lexer.LESS:
		return left.(float64) < right.(float64)
	case lexer.LESS_EQUAL:
		return left.(float64) <= right.(float64)

	}

	// should be unreachable
	return nil
}
