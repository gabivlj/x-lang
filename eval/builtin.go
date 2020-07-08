package eval

import (
	"xlang/object"
)

var builtins = map[string]*object.Builtin{
	"len": object.GetBuiltinByName("len"),

	"push": object.GetBuiltinByName("push"),

	"pop": object.GetBuiltinByName("pop"),

	"shift": object.GetBuiltinByName("shift"),

	"unshift": object.GetBuiltinByName("unshift"),

	"first": object.GetBuiltinByName("first"),

	"set": object.GetBuiltinByName("set"),

	"log": object.GetBuiltinByName("log"),

	"keys": object.GetBuiltinByName("keys"),

	"delete": object.GetBuiltinByName("delete"),
}
