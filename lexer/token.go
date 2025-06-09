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
	Line    int
}

func (t *Token) String() string {
	s := fmt.Sprintf("%v %s ", t.Type, t.Lexeme)
	s += PrintLiteral(t.Literal)

	return s
}

// helper function to stringify token's Literal value by converting to string based on type
func PrintLiteral(literal any) string {
	switch l := literal.(type) {
	case string:
		return l
	case float64:
		nStr := strconv.FormatFloat(l, 'f', -1, 64)
		// force decimal even when no decimal needed for representation to match java spec
		if !strings.Contains(nStr, ".") {
			nStr += ".0"
		}
		return nStr
	case bool:
		return strconv.FormatBool(l)
	case nil:
		return "null"
	default:
		e := fmt.Sprintf("Unknown value type in literal: %v", l)
		panic(e)
	}
}
