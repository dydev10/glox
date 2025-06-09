package interpreter

import (
	"fmt"

	"github.com/dydev10/glox/lexer"
)

type RuntimeError struct {
	token   *lexer.Token
	message string
}

func (re *RuntimeError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", re.token.Line, re.message)
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", re.token.Line, re.message)
}
