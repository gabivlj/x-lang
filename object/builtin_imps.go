package object

import "fmt"

// Len is the standard implementation of len(...) in Xlang
func Len(args ...Object) Object {

	if len(args) != 1 {
		return NewError("Expected 1 argument, got %d", len(args))
	}

	switch newObject := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(newObject.Value))}
	case *Array:
		return &Integer{Value: int64(len(newObject.Elements))}
	}
	return NewError("Unexpected type: %s for function len()", args[0].Type())
}

// Push is the standard imp. of push(...) in Xlang
func Push(args ...Object) Object {
	if len(args) < 1 {
		return NewError("Expected 1 argument or more, got %d", len(args))
	}
	switch newObject := args[0].(type) {
	case *Array:
		newElements := append(newObject.Elements[:], args[1:]...)
		return &Array{Elements: newElements}
	}
	return NewError("Unexpected type for push(); got %s", args[0].Type())
}

func array(o Object) (*Array, bool) {
	arr, ok := o.(*Array)
	if !ok {
		return nil, ok
	}
	return arr, ok
}

func function(o Object) (*Function, bool) {
	funct, ok := o.(*Function)
	if !ok {
		return nil, ok
	}
	return funct, ok
}

// First returns first arr element
func First(args ...Object) Object {
	if len(args) < 1 {
		return NewError("Expected 1 argument or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return NewError("Unexpected type for first(); got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return nil
	}
	return arr.Elements[0]
}

// Shift deletes first element
func Shift(args ...Object) Object {
	if len(args) < 1 {
		return NewError("Expected 1 argument or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return NewError("Unexpected type for shift(); got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return nil
	}
	return &Array{Elements: arr.Elements[1:]}
}

// Pop deletes last element of array
func Pop(args ...Object) Object {
	if len(args) < 1 {
		return NewError("Expected 1 argument or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		return NewError("Unexpected type for pop(); got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return arr
	}
	return &Array{Elements: arr.Elements[:len(arr.Elements)-1]}
}

// Unshift adds an element from the beginning of the array
func Unshift(args ...Object) Object {
	if len(args) < 2 {
		return NewError("Expected 2 arguments or more, got %d", len(args))
	}

	arr, ok := array(args[0])
	if !ok {
		return NewError("Unexpected type for unshift(); got %s", args[0].Type())
	}
	elements := append(args[1:], arr.Elements[:]...)
	return &Array{Elements: elements}
}

// Set sets a specific item of the array to the third parameter of the function
func Set(args ...Object) Object {
	if len(args) < 3 {
		return NewError("Expected 3 arguments or more, got %d", len(args))
	}
	arr, ok := array(args[0])
	if !ok {
		hash, k := args[0].(*HashMap)
		if k {
			return SetHash(hash, args[1], args[2])
		}
		return NewError("Unexpected type for first(); got %s", args[0].Type())
	}
	number, ok := args[1].(*Integer)
	if !ok {
		return NewError("Unexpected type for set(); got %s", args[0].Type())
	}
	if int(number.Value) >= len(arr.Elements) || int(number.Value) < 0 {
		return nil
	}
	arrElements := make([]Object, len(arr.Elements))
	copy(arrElements, arr.Elements)
	newArr := &Array{Elements: arrElements}

	newArr.Elements[number.Value] = args[2]
	return newArr
}

// SetHash creates a new entry in the hashmap
func SetHash(arg *HashMap, newKey Object, setVal Object) Object {
	hashable, ok := newKey.(Hashable)
	if !ok {
		arg.UnhashablePairs[newKey] = HashPair{Key: newKey, Value: setVal}
		return setVal
	}
	hashed := hashable.HashKey()
	arg.Pairs[hashed] = HashPair{Key: newKey, Value: setVal}
	return setVal
}

// Keys is the keys() function, accepts one hashtable and returns the keys of this
func Keys(args ...Object) Object {
	if len(args) != 1 {
		return NewError("Error: Expected 1 argument on keys() but got %d", len(args))
	}
	hashTable, ok := args[0].(*HashMap)
	if !ok {
		return NewError("Error: Expected object of type HashTable on keys(), got: %s", args[0].Type())
	}
	arr := Array{Elements: make([]Object, 0, len(hashTable.Pairs)+len(hashTable.UnhashablePairs))}
	for _, obj := range hashTable.Pairs {
		key := obj.Key
		arr.Elements = append(arr.Elements, key)
	}
	for _, obj := range hashTable.UnhashablePairs {
		key := obj.Key
		arr.Elements = append(arr.Elements, key)
	}
	return &arr
}

// Delete a key from a hash table
func Delete(args ...Object) Object {
	if len(args) != 2 {
		return NewError("Error: Expected 2 argument on delete() but got %d", len(args))
	}
	hash, ok := args[0].(*HashMap)
	if !ok {
		return NewError("Error: Expected HashMap as firs argument on delete() but got %s", args[0].Type())
	}
	hashable, ok := args[1].(Hashable)
	if ok {
		hashed := hashable.HashKey()
		h, ok := hash.Pairs[hashed]
		if !ok {
			return nil
		}
		delete(hash.Pairs, hashed)
		return h.Value
	}
	h, ok := hash.UnhashablePairs[args[1]]
	if !ok {
		return nil
	}
	delete(hash.UnhashablePairs, args[1])
	return h.Value
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// GetBuiltinByName returns the function builtin
func GetBuiltinByName(name string) *Builtin {
	for _, def := range builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

// Last returns the last element of an array
func Last(args ...Object) Object {
	if len(args) != 1 {
		return NewError("Error: Expected 1 argument on keys() but got %d", len(args))
	}
	arr, k := args[0].(*Array)
	if !k {
		return NewError("Error: Expected Array as first argument on last() but got %s", args[0].Type())
	}
	if len(arr.Elements) == 0 {
		return nil
	}
	return arr.Elements[len(arr.Elements)-1]
}
