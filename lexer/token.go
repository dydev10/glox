package lexer

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

func (t *Token) String() string {
	s := fmt.Sprintf("%v %s ", t.Type, t.Lexeme)
	if t.Literal != "" {
		s += t.Literal
	} else {
		s += "null"
	}
	return s
}
