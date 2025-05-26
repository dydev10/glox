package lexer

import "fmt"

type LexError struct {
	line    int
	message string
}

func (le *LexError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", le.line, le.message)
}
