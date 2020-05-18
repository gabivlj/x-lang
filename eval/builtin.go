package eval

import "xlang/object"

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
}
