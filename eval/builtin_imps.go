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
		hash, k := args[0].(*object.HashMap)
		if k {
			return SetHash(hash, args[1], args[2])
		}
		return object.NewError("Unexpected type for first(); got %s", args[0].Type())
	}
	number, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewError("Unexpected type for set(); got %s", args[0].Type())
	}
	if int(number.Value) >= len(arr.Elements) || int(number.Value) < 0 {
		return NULL
	}
	arrElements := make([]object.Object, len(arr.Elements))
	copy(arrElements, arr.Elements)
	newArr := &object.Array{Elements: arrElements}

	newArr.Elements[number.Value] = args[2]
	return newArr
}

// SetHash creates a new entry in the hashmap
func SetHash(arg *object.HashMap, newKey object.Object, setVal object.Object) object.Object {
	hashable, ok := newKey.(object.Hashable)
	if !ok {
		arg.UnhashablePairs[newKey] = object.HashPair{Key: newKey, Value: setVal}
		return setVal
	}
	hashed := hashable.HashKey()
	arg.Pairs[hashed] = object.HashPair{Key: newKey, Value: setVal}
	return setVal
}

// Keys is the keys() function, accepts one hashtable and returns the keys of this
func Keys(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Error: Expected 1 argument on keys() but got %d", len(args))
	}
	hashTable, ok := args[0].(*object.HashMap)
	if !ok {
		return object.NewError("Error: Expected object of type HashTable on keys(), got: %s", args[0].Type())
	}
	arr := object.Array{Elements: make([]object.Object, 0, len(hashTable.Pairs)+len(hashTable.UnhashablePairs))}
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
func Delete(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError("Error: Expected 2 argument on delete() but got %d", len(args))
	}
	hash, ok := args[0].(*object.HashMap)
	if !ok {
		return object.NewError("Error: Expected HashMap as firs argument on delete() but got %s", args[0].Type())
	}
	hashable, ok := args[1].(object.Hashable)
	if ok {
		hashed := hashable.HashKey()
		h, ok := hash.Pairs[hashed]
		if !ok {
			return NULL
		}
		delete(hash.Pairs, hashed)
		return h.Value
	}
	h, ok := hash.UnhashablePairs[args[1]]
	if !ok {
		return NULL
	}
	delete(hash.UnhashablePairs, args[1])
	return h.Value
}
