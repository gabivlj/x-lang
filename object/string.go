package object

// String is a string object
type String struct {
	Value string
}

// Type .
func (s *String) Type() ObjectType {
	return StringObject
}

// Inspect .
func (s *String) Inspect() string { 
	return s.Value
}