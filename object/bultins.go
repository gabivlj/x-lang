package object

var builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: Len},
	},

	{
		"push",
		&Builtin{Fn: Push},
	},

	{"pop",
		&Builtin{Fn: Pop},
	},
	{"shift",
		&Builtin{Fn: Shift},
	},
	{"unshift",
		&Builtin{Fn: Unshift},
	},
	{"first",
		&Builtin{Fn: First},
	},
	{"set",
		&Builtin{Fn: Set},
	},
	{
		"last",
		&Builtin{Fn: Last},
	},
	{
		"log",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					arg.Inspect()
				}
				if len(args) == 0 {
					return &Log{Message: nil}
				}
				if len(args) == 1 {
					return &Log{Message: args[0]}
				}
				arr := Array{Elements: make([]Object, 0, len(args))}
				for _, arg := range args {
					arr.Elements = append(arr.Elements, arg)
				}
				return &Log{Message: &arr}
			},
		},
	},

	{"keys",
		&Builtin{Fn: Keys},
	},

	{"delete",
		&Builtin{Fn: Delete},
	},
}

// GetBuiltins objects
func GetBuiltins() []struct {
	Name    string
	Builtin *Builtin
} {
	return builtins
}
