package eval

import (
	"xlang/object"
)

// Len is the standard implementation of len(...) in Xlang
func Len(args ...object.Object) object.Object {

	if len(args) != 1 {
		return object.NewError("Expected 1 argument, got %d", len(args))
	}

	switch newObject := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(newObject.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(newObject.Elements))}
	}
	return object.NewError("Unexpected type: %s for function len()", args[0].Type())
}

// Push is the standard imp. of push(...) in Xlang
func Push(args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("Expected 1 argument or more, got %d", len(args))
	}
	switch newObject := args[0].(type) {
	case *object.Array:
		newElements := append(newObject.Elements[:], args[1:]...)
		return &object.Array{Elements: newElements}
	}
	return object.NewError("Unexpected type for push(); got %s", args[0].Type())
}

func array(o object.Object) (*object.Array, bool) {
	arr, ok := o.(*object.Array)
	if !ok {
		return nil, ok
	}
	return arr, ok
}

func function(o object.Object) (*object.Function, bool) {
	funct, ok := o.(*object.Function)
	if !ok {
		return nil, ok
	}
	return funct, ok
}

// First returns first arr element
func First(args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("Expected 1 argument or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return object.NewError("Unexpected type for first(); got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return NULL
	}
	return arr.Elements[0]
}

// Shift deletes first element
func Shift(args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("Expected 1 argument or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return object.NewError("Unexpected type for shift(); got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return arr
	}
	return &object.Array{Elements: arr.Elements[1:]}
}

// Pop deletes last element of array
func Pop(args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("Expected 1 argument or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return object.NewError("Unexpected type for pop(); got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return arr
	}
	return &object.Array{Elements: arr.Elements[:len(arr.Elements)-1]}
}

// Unshift adds an element from the beginning of the array
func Unshift(args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError("Expected 2 arguments or more, got %d", len(args))
	}

	arr, ok := array(args[0])
	if !ok {
		return object.NewError("Unexpected type for unshift(); got %s", args[0].Type())
	}
	elements := append(args[1:], arr.Elements[:]...)
	return &object.Array{Elements: elements}
}

// Set sets a specific item of the array to the third parameter of the function
func Set(args ...object.Object) object.Object {
	if len(args) < 3 {
		return object.NewError("Expected 3 arguments or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return object.NewError("Unexpected type for first(); got %s", args[0].Type())
	}
	number, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewError("Unexpected type for set(); got %s", args[0].Type())
	}
	if int(number.Value) >= len(arr.Elements) || int(number.Value) < 0 {
		return NULL
	}
	arr.Elements[number.Value] = args[2]
	return arr.Elements[number.Value]
}
