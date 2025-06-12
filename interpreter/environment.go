package interpreter

import (
	"fmt"

	"github.com/dydev10/glox/lexer"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (env *Environment) define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) get(name *lexer.Token) (any, error) {
	if val, ok := env.values[name.Lexeme]; ok {
		return val, nil
	} else {
		return nil, &RuntimeError{
			token:   name,
			message: fmt.Sprintf("Undefined variable %s.", name.Lexeme),
		}
	}
}

func (env *Environment) assign(name *lexer.Token, value any) error {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return nil
	}

	return &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable %s.", name.Lexeme),
	}
}
