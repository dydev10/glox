package interpreter

import (
	"fmt"

	"github.com/dydev10/glox/lexer"
)

type ResolveError struct {
	token   *lexer.Token
	message string
}

func (re *ResolveError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", re.token.Line, re.message)
}

func (re *ResolveError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", re.token.Line, re.message)
}
