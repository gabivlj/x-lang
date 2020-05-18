package object

// BuiltinFunction is built in code inside Xlang
type BuiltinFunction func(args ...Object) Object

// Builtin is a builtin function in Xlang
type Builtin struct {
	Fn BuiltinFunction
}

// Type .
func (b *Builtin) Type() ObjectType { return BuiltinObject }

// Inspect .
func (b *Builtin) Inspect() string { return "builtin function" }
