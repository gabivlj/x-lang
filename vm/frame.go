package vm

import (
	"xlang/code"
	"xlang/object"
)

// Frame stores the function thats is being called information
type Frame struct {
	fn *object.Closure
	// Current position in the bytecode
	ip int
	// Stores where the function is stored in the stack
	basePointer int
}

// NewFrame ...
func NewFrame(fn *object.Closure, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

// Instructions .
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Fn.Instructions
}
