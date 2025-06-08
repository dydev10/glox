package parser

import (
	"fmt"

	"github.com/dydev10/glox/lexer"
)

type ParseError struct {
	token   *lexer.Token
	message string
}

func (le *ParseError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", le.token.Line, le.message)
}

func (le *ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", le.token.Line, le.message)
}
