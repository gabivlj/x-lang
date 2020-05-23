package eval

import (
	"fmt"
	"xlang/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: Len,
	},

	"push": {
		Fn: Push,
	},

	// TODO: Implement first, last, pop, unshift, shift, for_each, for, map, reduce
	"pop": {
		Fn: Pop,
	},
	"shift": {
		Fn: Shift,
	},
	"unshift": {
		Fn: Unshift,
	},
	"first": {
		Fn: First,
	},
	"set": {
		Fn: Set,
	},
	"log": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println("QIE ES", arg.Inspect())
			}
			if len(args) == 0 {
				return &object.Log{Message: NULL}
			}
			if len(args) == 1 {
				return &object.Log{Message: args[0]}
			}
			arr := object.Array{Elements: make([]object.Object, 0, len(args))}
			for _, arg := range args {
				arr.Elements = append(arr.Elements, arg)
			}
			return &object.Log{Message: &arr}
		},
	},

	"keys": {
		Fn: Keys,
	},

	"delete": {
		Fn: Delete,
	},
}
