package interpreter

import (
	"fmt"

	"github.com/dydev10/glox/lexer"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

func (env *Environment) define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) get(name *lexer.Token) (any, error) {
	if val, ok := env.values[name.Lexeme]; ok {
		return val, nil
	}
	if env.enclosing != nil {
		return env.enclosing.get(name)
	}
	return nil, &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable %s.", name.Lexeme),
	}
}

func (env *Environment) assign(name *lexer.Token, value any) error {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.assign(name, value)
	}
	return &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable %s.", name.Lexeme),
	}
}
