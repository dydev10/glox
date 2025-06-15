package interpreter

type LoxCallable interface {
	Arity() int
	Call(intr *Interpreter, arguments []any) (any, error)
}
