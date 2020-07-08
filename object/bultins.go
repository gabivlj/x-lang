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
	//  {
	// 	"log",
	// 	&Builtin{
	// 		Fn: func(args ...object.Object) object.Object {
	// 			for _, arg := range args {
	// 				arg.Inspect())fmt.Println("QIE ES",
	// 			}
	// 			if len(args) == 0 {
	// 				return &object.Log{Message: NULL}
	// 			}
	// 			if len(args) == 1 {
	// 				return &object.Log{Message: args[0]}
	// 			}
	// 			arr := object.Array{Elements: make([]object.Object, 0, len(args))}
	// 			for _, arg := range args {
	// 				arr.Elements = append(arr.Elements, arg)
	// 			}
	// 			return &object.Log{Message: &arr}
	// 		}
	// 	},
	// },

	{"keys",
		&Builtin{Fn: Keys},
	},

	{"delete",
		&Builtin{Fn: Delete},
	},
}
