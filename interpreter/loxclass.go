package interpreter

type LoxClass struct {
	name string
}

func (c *LoxClass) Arity() int {
	// TODO: implement this method
	return 0
}

func (c *LoxClass) Call(intr *Interpreter, arguments []any) (any, error) {
	// TODO: implement this method
	return nil, nil
}

func (c *LoxClass) String() string {
	return c.name
}
