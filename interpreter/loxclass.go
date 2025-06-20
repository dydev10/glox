package interpreter

type LoxClass struct {
	name string
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) Call(intr *Interpreter, arguments []any) (any, error) {
	instance := &LoxInstance{
		class:  c,
		fields: make(map[string]any),
	}

	return instance, nil
}

func (c *LoxClass) String() string {
	return c.name
}
