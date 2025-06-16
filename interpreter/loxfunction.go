package interpreter

import "github.com/dydev10/glox/ast"

type LoxFunction struct {
	declaration *ast.Function
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(intr *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironment(intr.globals)

	for i := range f.declaration.Params {
		environment.define(f.declaration.Params[i].Lexeme, arguments[i])
	}

	if err := intr.executeBlock(f.declaration.Body, environment); err != nil {
		return nil, err
	}

	return nil, nil
}

func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
