package object

import "strings"

// Array stores inside an array objects
type Array struct {
	Elements []Object
}

// Type .
func (a *Array) Type() ObjectType { return ArrayObject }

// Inspect .
func (a *Array) Inspect() string {
	var builder strings.Builder
	builder.WriteByte('[')
	for idx, obj := range a.Elements {
		builder.WriteString(obj.Inspect())
		if idx == len(a.Elements)-1 {
			break
		}
		builder.WriteByte(',')
	}
	builder.WriteByte(']')
	return builder.String()
}
