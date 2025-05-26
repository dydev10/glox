package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

func (t *Token) String() string {
	s := fmt.Sprintf("%v %s ", t.Type, t.Lexeme)

	switch l := t.Literal.(type) {
	case string:
		s += l
	case float64:
		nStr := strconv.FormatFloat(l, 'f', -1, 64)
		// force decimal even when no decimal needed for representation to match java spec
		if !strings.Contains(nStr, ".") {
			nStr += ".0"
		}
		s += nStr
	case nil:
		s += "null"
	default:
		panic("Unknown value type in literal")
	}

	return s
}
