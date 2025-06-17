package interpreter

type ThrownReturn struct {
	value any
}

func (re *ThrownReturn) Error() string {
	return "Not an error. If this shows up as error in logs, something is wrong with handling of return value"
}
