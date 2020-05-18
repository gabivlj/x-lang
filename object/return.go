package object

// ReturnValue is the value of a return
type ReturnValue struct {
	Value Object
}

// Type .
func (rv *ReturnValue) Type() ObjectType {
	return ReturnObject
}

// Inspect .
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
