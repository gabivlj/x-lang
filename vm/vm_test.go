package vm

import (
	"fmt"
	"testing"
	"xlang/ast"
	"xlang/compiler"
	"xlang/lexer"
	"xlang/object"
	"xlang/parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVMTests(t *testing.B, tests []vmTestCase, printStackTraceAndStop ...bool) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compile error: %s", err.Error())
		}
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}
		stackElem := vm.LastPoppedStackElem()
		if stackElem == nil {
			t.Fatalf("error, stackElement is null, expected %#v in %#v and ins %s", tt.expected, tt.input, vm.instructions.String())
		}
		testExpectedObject(t, tt.expected, stackElem)
		if len(printStackTraceAndStop) > 0 && t.Failed() && printStackTraceAndStop[0] {
			t.Fatal("Failed on: ", tt.input)
		}
	}
}

func testExpectedObject(
	t *testing.B,
	expected interface{},
	actual object.Object,
) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)

		}
	}
}

func testBooleanObject(expected bool, actual object.Object) error {
	if actual.Type() != object.BooleanObject {
		return fmt.Errorf("expected type=%s. got=%s", object.BooleanObject, actual.Type())
	}
	if expected != actual.(*object.Boolean).Value {
		return fmt.Errorf("expected value=%t. got=%t", expected, !expected)
	}
	return nil
}

func BenchmarkBooleanExpressions(t *testing.B) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", false},
	}

	runVMTests(t, tests)
}

func BenchmarkIntegerArithmetic(t *testing.B) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
	}
	runVMTests(t, tests)
}
