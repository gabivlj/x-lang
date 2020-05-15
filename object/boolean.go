package object

import (
	"fmt"
)

// Boolean type
type Boolean struct {
	Value bool
}

// Type boolean
func (b *Boolean) Type() ObjectType {
	return BooleanObject
}

// Inspect boolean
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
