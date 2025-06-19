package interpreter

import (
	"fmt"
	"os"

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

func (env *Environment) ancestor(distance int) *Environment {
	environment := env
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
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

func (env *Environment) getAt(distance int, name string) (any, error) {
	value, ok := env.ancestor(distance).values[name]
	if !ok {
		// TODO: accept token instead of token's name and return runtime error here when variable not found in environment
		fmt.Fprintf(os.Stderr, "Unhandled Error: getAt function did not find variable %s definition in given depth %d", name, distance)
	}
	return value, nil
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

func (env *Environment) assignAt(distance int, name *lexer.Token, value any) {
	env.ancestor(distance).values[name.Lexeme] = value
}
