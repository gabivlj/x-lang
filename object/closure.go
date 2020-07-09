package object

import "fmt"

// Closure stores the free variables of a function and the compiled function
type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

// Type .
func (c *Closure) Type() ObjectType { return ClosureObject }

// Inspect .
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}
