package interpreter

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
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

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}

	return nil
}
