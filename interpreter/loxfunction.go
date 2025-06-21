package interpreter

import "github.com/dydev10/glox/ast"

type LoxFunction struct {
	declaration   *ast.Function
	closure       *Environment
	isInitializer bool
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(intr *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironment(f.closure)

	for i := range f.declaration.Params {
		environment.define(f.declaration.Params[i].Lexeme, arguments[i])
	}

	if err := intr.executeBlock(f.declaration.Body, environment); err != nil {
		thrownReturn, isReturn := err.(*ThrownReturn)
		if isReturn {
			// if its a constructor, ignore return value and just return 'this'. return value syntax should be block by resolver
			if f.isInitializer {
				return f.closure.getAt(0, "this")
			}
			return thrownReturn.value, nil
		}
		return nil, err
	}

	// always return 'this' if its constructor
	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}

	return nil, nil
}

func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.closure)
	environment.define("this", instance)

	return &LoxFunction{
		declaration:   f.declaration,
		closure:       environment,
		isInitializer: f.isInitializer,
	}
}
