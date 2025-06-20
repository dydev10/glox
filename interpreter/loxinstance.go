package interpreter

import (
	"fmt"

	"github.com/dydev10/glox/lexer"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func (i *LoxInstance) Get(name *lexer.Token) (any, error) {
	if val, ok := i.fields[name.Lexeme]; ok {
		return val, nil
	}

	return nil, &RuntimeError{
		token:   name,
		message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme),
	}
}

func (i *LoxInstance) Set(name *lexer.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) String() string {
	return i.class.name + " instance"
}
