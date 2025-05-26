package lexer

import "fmt"

type LexError struct {
	line int
	ch   rune
}

func (le *LexError) String() string {
	return fmt.Sprintf("[line %d] Error: Unexpected character: %s", le.line, string(le.ch))
}
