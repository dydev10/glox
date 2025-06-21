package interpreter

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func (c *LoxClass) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (c *LoxClass) Call(intr *Interpreter, arguments []any) (any, error) {
	instance := &LoxInstance{
		class:  c,
		fields: make(map[string]any),
	}
	// call constructor method of class after binding 'this'
	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(intr, arguments)
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
