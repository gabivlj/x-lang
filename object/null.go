package object

// Null .
type Null struct{}

// Type .
func (n *Null) Type() ObjectType { return NullObject }

// Inspect null
func (n *Null) Inspect() string { return "null" }
