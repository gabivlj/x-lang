package vm

import (
	"encoding/binary"
	"fmt"
	"xlang/code"
	"xlang/compiler"
	"xlang/object"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

// VM holds all the Virtual Machine information and logic
type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	// Points to the next value. Top of stack is stack[sp-1]
	sp int
}

// New returns a new VM from a bytecode
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// StackTop returns the top element on the stack
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// Run runs the VM
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			{
				ip++
				idx := binary.BigEndian.Uint16(vm.instructions[ip:])
				ip++
				if err := vm.push(vm.constants[idx]); err != nil {
					return err
				}
			}
		case code.OpGreaterThan:
			{
				left, right, err := vm.unwrapTwoIntegers()
				if err != nil {
					return err
				}
				if err := vm.push(nativeToBooleanObject(left.Value > right.Value)); err != nil {
					return err
				}
			}
		case code.OpNotEqual, code.OpEqual:
			{
				right := vm.pop()
				left := vm.pop()
				if leftInteger, k := left.(*object.Integer); k {
					if err := vm.numericalComparison(leftInteger, right, op); err != nil {
						return err
					}
					continue
				}
				// Standard comparison
				equal := right == left
				if code.OpNotEqual == op {
					equal = !equal
				}
				if err := vm.push(nativeToBooleanObject(equal)); err != nil {
					return err
				}
			}
		case code.OpTrue:
			{
				if err := vm.push(True); err != nil {
					return err
				}
			}
		case code.OpFalse:
			{
				if err := vm.push(False); err != nil {
					return err
				}
			}
		case code.OpPop:
			{
				vm.pop()
			}
		case code.OpDiv:
			{
				left, right, err := vm.unwrapTwoIntegers()
				if err != nil {
					return err
				}
				val := left.Value / right.Value
				if err := vm.push(&object.Integer{Value: val}); err != nil {
					return err
				}
			}
		case code.OpMul:
			{
				left, right, err := vm.unwrapTwoIntegers()
				if err != nil {
					return err
				}
				val := left.Value * right.Value
				if err := vm.push(&object.Integer{Value: val}); err != nil {
					return err
				}
			}
		case code.OpSub:
			{
				left, right, err := vm.unwrapTwoIntegers()
				if err != nil {
					return err
				}
				val := left.Value - right.Value
				if err := vm.push(&object.Integer{Value: val}); err != nil {
					return err
				}
			}
		case code.OpAdd:
			{
				right := vm.pop()
				left := vm.pop()
				switch rightObject := right.(type) {
				case *object.Integer:
					{
						leftObject, ok := left.(*object.Integer)
						if !ok {
							return fmt.Errorf("expected %s, got: %s", right.Type(), left.Type())
						}
						if err := vm.push(&object.Integer{Value: leftObject.Value + rightObject.Value}); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (vm *VM) numericalComparison(leftInteger *object.Integer, right object.Object, op code.Opcode) error {
	rightInteger, k := right.(*object.Integer)
	if !k {
		return fmt.Errorf("expected numerical value, got=%s", rightInteger.Type())
	}
	equal := leftInteger.Value == rightInteger.Value
	if op == code.OpNotEqual {
		equal = !equal
	}
	return vm.push(nativeToBooleanObject(equal))
}

func nativeToBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}

func (vm *VM) popIntegerObject() (*object.Integer, error) {
	obj := vm.pop()
	parsedInteger, k := obj.(*object.Integer)
	if !k {
		return nil, fmt.Errorf("expected integer object, got=%s", obj.Type())
	}
	return parsedInteger, nil
}

func (vm *VM) unwrapTwoIntegers() (*object.Integer, *object.Integer, error) {
	right, err := vm.popIntegerObject()
	if err != nil {
		return nil, nil, err
	}
	left, err := vm.popIntegerObject()
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

// LastPoppedStackElem returns the last popped element
func (vm *VM) LastPoppedStackElem() object.Object {
	// for _, obj := range vm.stack {
	// 	fmt.Println("res", obj.Inspect(), vm.sp)
	// }
	return vm.stack[vm.sp]
}
