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
		// for i, constant := range comp.Bytecode().Constants {
		// 	fmt.Printf("CONSTANT %d %p (%T):\n", i, constant, constant)

		// 	switch constant := constant.(type) {
		// 	case *object.CompiledFunction:
		// 		fmt.Printf(" Instructions:\n%s", constant.Instructions)
		// 	case *object.Integer:
		// 		fmt.Printf(" Value: %d\n", constant.Value)
		// 	}

		// 	fmt.Printf("\n")
		// }
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Errorf("vm error: %s", err)
			if len(printStackTraceAndStop) > 0 && t.Failed() && printStackTraceAndStop[0] {
				t.Fatal("Failed on: ", tt.input)
			}
		}
		stackElem := vm.LastPoppedStackElem()
		if stackElem == nil {
			t.Fatalf("error, stackElement is null, expected %#v in %#v and ins %s", tt.expected, tt.input, vm.currentFrame().Instructions().String())
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
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}

	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not Error: %T (%+v)", actual, actual)
			return
		}
		if errObj.Message != expected.Message {
			t.Errorf("wrong error message. expected=%q, got=%q",
				expected.Message, errObj.Message)
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.HashMap)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of Pairs. want=%d, got=%d",
				len(expected), len(hash.Pairs))
			return
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}

			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements. want=%d, got=%d",
				len(expected), len(array.Elements))
			return
		}

		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	}
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
	}

	return nil
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
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runVMTests(t, tests, true)
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
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVMTests(t, tests)
}

func BenchmarkConditionals(t *testing.B) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVMTests(t, tests)
}

func BenchmarkGlobalLetStatements(t *testing.B) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = if (false) { 1 } else { 2 }; let two = one + one; one + two", 6},
	}

	runVMTests(t, tests)
}

func BenchmarkStringExpressions(t *testing.B) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVMTests(t, tests)
}

func BenchmarkArrayLiterals(t *testing.B) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVMTests(t, tests)
}

func BenchmarkIndexExpressions(t *testing.B) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVMTests(t, tests, true)
}

func BenchmarkHashLiterals(t *testing.B) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1: 2, 2: 3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVMTests(t, tests)
}

func BenchmarkCallingFunctionsWithoutArguments(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
	let earlyExit = fn() { return 99; 100; };
	earlyExit();
	`,
			expected: 99,
		},
		{
			input: `
	let earlyExit = fn() { return 99; return 100; };
	earlyExit();
	`,
			expected: 99,
		},
		{
			input: `
	let one = fn() { 1; };
	let two = fn() { 2; };
	one() + two()
	`,
			expected: 3,
		},
		{
			input: `
	let a = fn() { 1 };
	let b = fn() { a() + 1 };
	let c = fn() { b() + 1 };
	c();
	`,
			expected: 3,
		},
		{
			input: `
			let fivePlusTen = fn() { 5 + 10; };
			fivePlusTen();
			`,
			expected: 15,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkFunctionsWithoutReturnValue(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn = fn() { };
			noReturn();
			`,
			expected: Null,
		},
		{
			input: `
			let noReturn = fn() { };
			let noReturnTwo = fn() { noReturn(); };
			noReturn();
			noReturnTwo();
			`,
			expected: Null,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkFirstClassFunctions(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOne = fn() { 1; };
			let returnsOneReturner = fn() { returnsOne; };
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkCallingFunctionsWithBindings(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() { let one = 1; one };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
			oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() {
					let num = 1;
					globalSeed - num;
			}
			let minusTwo = fn() {
					let num = 2;
					globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkFirstClassFunctionsV2(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOneReturner = fn() {
					let returnsOne = fn() { 1; };
					returnsOne;
			};
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkCallingFunctionsWithArgumentsAndBindings(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let identity = fn(a) { a; };
			identity(4);
			`,
			expected: 4,
		},
		{
			input: `
			let sum = fn(a, b) { a + b; };
			sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
	let sum = fn(a, b) {
			let c = a + b;
			c;
	};
	sum(1, 2);
	`,
			expected: 3,
		},
		{
			input: `
	let sum = fn(a, b) {
			let c = a + b;
			c;
	};
	sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input: `
	let sum = fn(a, b) {
			let c = a + b;
			c;
	};
	let outer = fn() {
			sum(1, 2) + sum(3, 4);
	};
	outer();
	`,
			expected: 10,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkCallingFunctionsWithWrongArguments(t *testing.B) {
	tests := []vmTestCase{
		{
			input:    `fn() { 1; }(1);`,
			expected: `wrong number of parameters, expected=0, got=1`,
		},
		{
			input:    `fn(a) { a; }();`,
			expected: `wrong number of parameters, expected=1, got=0`,
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: `wrong number of parameters, expected=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func BenchmarkBuiltinFunctions(t *testing.B) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{
			`len(1)`,
			&object.Error{
				Message: "Unexpected type: INTEGER for function len()",
			},
		},
		{`len("one", "two")`,
			&object.Error{
				Message: "Expected 1 argument, got 2",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`push(["hello"], "world!")`, &object.Array{Elements: []object.Object{&object.String{Value: "hello"}, &object.String{Value: "world!"}}}},
		{`first([1, 2, 3])`, 1},
		{`first([])`, Null},
		{`first(1)`,
			&object.Error{
				Message: "Unexpected type for first(); got INTEGER",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`, Null},
		{`last(1)`,
			&object.Error{
				Message: "Error: Expected Array as first argument on last() but got INTEGER",
			},
		},
		{`shift([1, 2, 3])`, []int{2, 3}},
		{`shift([])`, Null},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`,
			&object.Error{
				Message: "Unexpected type for push(); got INTEGER",
			},
		},
	}

	runVMTests(t, tests, true)
}

func BenchmarkClosures(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let newClosure = fn(a) {
					fn() { a; };
			};
			let closure = newClosure(99);
			closure();
			`,
			expected: 99,
		},
		{
			input: `
	let newAdder = fn(a, b) {
			fn(c) { a + b + c };
	};
	let adder = newAdder(1, 2);
	adder(8);
	`,
			expected: 11,
		},
		{
			input: `
	let newAdder = fn(a, b) {
			let c = a + b;
			fn(d) { c + d };
	};
	let adder = newAdder(1, 2);
	adder(8);
	`,
			expected: 11,
		},
		{
			input: `
	let newAdderOuter = fn(a, b) {
			let c = a + b;
			fn(d) {
					let e = d + c;
					fn(f) { e + f; };
			};
	};
	let newAdderInner = newAdderOuter(1, 2)
	let adder = newAdderInner(3);
	adder(8);
	`,
			expected: 14,
		},
		{
			input: `
	let a = 1;
	let newAdderOuter = fn(b) {
			fn(c) {
					fn(d) { a + b + c + d };
			};
	};
	let newAdderInner = newAdderOuter(2)
	let adder = newAdderInner(3);
	adder(8);
	`,
			expected: 14,
		},
		{
			input: `
	let newClosure = fn(a, b) {
			let one = fn() { a; };
			let two = fn() { b; };
			fn() { one() + two(); };
	};
	let closure = newClosure(9, 90);
	closure();
	`,
			expected: 99,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkRecursiveFunctions(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
				let countDown = fn(x) {
						if (x == 0) {
								return 0;
						} else {
								countDown(x - 1);
						}
				};
				countDown(1);
				`,
			expected: 0,
		},
		{
			input: `
		let countDown = fn(x) {
				if (x == 0) {
						return 0;
				} else {
						countDown(x - 1);
				}
		};
		let wrapper = fn() {
				countDown(1);
		};
		wrapper();
		`,
			expected: 0,
		},
		{
			input: `
	let wrapper = fn() {
			let countDown = fn(x) {
					if (x == 0) {
							return 2;
					} else {
							countDown(x - 1);
					}
			};
			countDown(1);
	};
	wrapper();
	`,
			expected: 2,
		},
	}

	runVMTests(t, tests)
}

func BenchmarkRecursiveFibonacci(t *testing.B) {
	tests := []vmTestCase{
		{
			input: `
			let fibonacci = fn(x) {
					if (x == 0) {
							return 0;
					} else {
							if (x == 1) {
									return 1;
							} else {
									return fibonacci(x - 1) + fibonacci(x - 2);
							}
					}
			};
			fibonacci(15);
			`,
			expected: 610,
		},
	}

	runVMTests(t, tests)
}
