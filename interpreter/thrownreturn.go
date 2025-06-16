package interpreter

import (
	"fmt"
)

type ThrownReturn struct {
	value any
}

func (re *ThrownReturn) Error() string {
	return fmt.Sprintf("Not an error. If this shows up as error in logs, something is wrong with handling of return value")
}
