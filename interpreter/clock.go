package interpreter

import "time"

type Clock struct {
}

func (clock *Clock) Arity() int {
	return 0
}

func (clock *Clock) Call(intr *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().Unix()), nil
}

func (clock *Clock) String() string {
	return "<native fn>"
}
