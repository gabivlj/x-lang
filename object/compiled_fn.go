package object

import (
	"fmt"
	"xlang/code"
)

// CompiledFunction stores Vm instructions
type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
}

// Type .
func (cf *CompiledFunction) Type() ObjectType { return CompiledFunctionObject }

// Inspect .
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
