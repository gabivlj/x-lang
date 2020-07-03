package vm

import (
	"encoding/binary"
	"fmt"
	"xlang/code"
	"xlang/compiler"
	"xlang/object"
)

const StackSize = 2048

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
				err := vm.push(vm.constants[idx])
				if err != nil {
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
						vm.push(&object.Integer{Value: leftObject.Value + rightObject.Value})
					}
				}
			}
		}
	}
	return nil
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
