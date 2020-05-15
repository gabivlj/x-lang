package object

import (
	"fmt"
)

// Integer representation in the xlang language
type Integer struct {
	Value int64
}

// Type returns the integer type
func (i *Integer) Type() ObjectType { return IntegerObject }

// Inspect inspects the value of integer
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
