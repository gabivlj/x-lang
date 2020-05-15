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
)

// Object is a xlang object.
type Object interface {
	Type() ObjectType
	Inspect() string
}
