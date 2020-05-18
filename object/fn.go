package object

import (
	"fmt"
	"xlang/ast"
)

// Function is an object which stores a function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns interface type
func (f *Function) Type() ObjectType {
	return FunctionObject
}

// Inspect inspects the function
func (f *Function) Inspect() string {
	str := ""
	if len(f.Parameters) == 0 {
		return fmt.Sprintf("fn () { \n %s \n}", f.Body.String())
	}
	for _, param := range f.Parameters {
		str += param.Value + ","
	}
	str = str[:len(str)-1]
	s := fmt.Sprintf("fn (%s) { \n %s \n}", str, f.Body.String())
	return s
}
