package object

// ObjectType is the string of the type
type ObjectType string

const (
	// IntegerObject is the integer type
	IntegerObject = "INTEGER"
	// NullObject is the null type
	NullObject = "OBJECT"
	// BooleanObject is the boolean type
	BooleanObject = "BOOL"
	// ReturnObject is the value wrapped around a return
	ReturnObject = "RETURN_VALUE"
	// ErrorObject is an error in running the AST
	ErrorObject = "ERROR"
	// FunctionObject is a function
	FunctionObject = "FUNCTION"
	// StringObject is a string representation in Xlang
	StringObject = "STRING"
	// BuiltinObject are things that are already implemented in the language
	BuiltinObject = "BUILTIN"

	// LogObject is an object which if you return with log(...) will log it into the logger.
	LogObject = "LOG"
	// ArrayObject is the builtin array system in Xlang
	ArrayObject = "ARRAY"
	// HashObject is a hashmap
	HashObject = "HASH"
	// CompiledFunctionObject is a function that stores instructions for the vm
	CompiledFunctionObject = "COMPILED FUNCTION"
	// ClosureObject is a function that stores a function and the freevariables
	ClosureObject = "CLOSURE"
)

// Object is a xlang object.
type Object interface {
	Type() ObjectType
	Inspect() string
}
