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
	BuiltinObject        = "BUILTIN"
	BuiltinObjectSpecial = "BUILTIN_SPECIAL"
	// ArrayObject is the builtin array system in Xlang
	ArrayObject = "ARRAY"
)

// Object is a xlang object.
type Object interface {
	Type() ObjectType
	Inspect() string
}
