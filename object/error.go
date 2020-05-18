package object

import "fmt"

// Error represents an error running the AST
type Error struct {
	Message string
}

// Type .
func (e *Error) Type() ObjectType { return ErrorObject }

// Inspect .
func (e *Error) Inspect() string { return fmt.Sprintf("Error: %s", e.Message) }

// NewError returns a new error
func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// IsError return if its an error
func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ErrorObject
	}
	return false
}
