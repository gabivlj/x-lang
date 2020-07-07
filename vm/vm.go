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
var Null = &object.Null{}

// VM holds all the Virtual Machine information and logic
type VM struct {
	constants []object.Object

	stack []object.Object
	// Points to the next value. Top of stack is stack[sp-1]
	sp          int
	globals     []object.Object
	frames      []*Frame
	framesIndex int
}

const GlobalsSize = 65536
const MaxFrames = 65536

// New returns a new VM from a bytecode
func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	return &VM{
		constants:   bytecode.Constants,
		stack:       make([]object.Object, StackSize),
		sp:          0,
		globals:     make([]object.Object, GlobalsSize),
		frames:      frames,
		framesIndex: 1,
	}
}

// NewWithGlobalsStore .
func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

// StackTop returns the top element on the stack
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func isTruthy(o object.Object) bool {
	switch o := o.(type) {
	case *object.Boolean:
		return o.Value
	case *object.Null:
		return false
	}
	return true
}

// Run runs the VM
func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++
		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])
		switch op {
		case code.OpIndex:
			{
				index := vm.pop()
				element := vm.pop()
				switch element := element.(type) {
				case *object.Array:
					{
						integerObject, ok := index.(*object.Integer)
						if !ok {
							return fmt.Errorf("expected integer got=%s", index.Type())
						}
						var objectToPush object.Object
						if integerObject.Value < 0 || integerObject.Value >= int64(len(element.Elements)) {
							objectToPush = Null
						} else {
							objectToPush = element.Elements[integerObject.Value]
						}
						if err := vm.push(objectToPush); err != nil {
							return err
						}
					}
				case *object.HashMap:
					{
						hashable, ok := index.(object.Hashable)
						var key object.HashKey
						if !ok {
							key = (&object.String{Value: index.Inspect()}).HashKey()
						} else {
							key = hashable.HashKey()
						}
						pair, ok := element.Pairs[key]
						if !ok {
							if err := vm.push(Null); err != nil {
								return err
							}
							continue
						}
						if err := vm.push(pair.Value); err != nil {
							return err
						}
					}
				default:
					return fmt.Errorf("invalid index operation on %s", element.Type())
				}
			}
		case code.OpHash:
			{
				lenOfHash := int(binary.BigEndian.Uint16(ins[ip+1:]))
				vm.currentFrame().ip += 2
				elements := vm.stack[vm.sp-lenOfHash : vm.sp]
				hash := make(map[object.HashKey]object.HashPair, len(elements)/2)
				for i := 0; i < len(elements); i += 2 {
					switch key := elements[i].(type) {
					case object.Hashable:
						hashed := key.HashKey()
						hash[hashed] = object.HashPair{Value: elements[i+1], Key: elements[i]}
					default:
						str := key.Inspect()
						obj := (&object.String{Value: str}).HashKey()
						hash[obj] = object.HashPair{Value: elements[i+1], Key: elements[i]}
					}
				}
				if err := vm.push(&object.HashMap{Pairs: hash}); err != nil {
					return err
				}
			}
		case code.OpArray:
			{
				lenOfArray := int(binary.BigEndian.Uint16(ins[ip+1:]))
				vm.currentFrame().ip += 2
				elements := vm.stack[vm.sp-lenOfArray : vm.sp]
				if err := vm.push(&object.Array{Elements: elements}); err != nil {
					return err
				}
			}
		case code.OpGetGlobal:
			{
				pos := int(binary.BigEndian.Uint16(ins[ip+1:]))
				vm.currentFrame().ip += 2
				obj := vm.globals[pos]
				if err := vm.push(obj); err != nil {
					return err
				}
			}

		case code.OpReturn:
			{
				frame := vm.popFrame()
				// Go to the starting point
				vm.sp = frame.basePointer - 1
				if err := vm.push(Null); err != nil {
					return err
				}
			}

		case code.OpReturnValue:
			{
				// Pop return value
				returnValue := vm.pop()
				// Pop the current frame (scope)
				frame := vm.popFrame()
				// Go to before we called the function
				vm.sp = frame.basePointer - 1
				err := vm.push(returnValue)
				if err != nil {
					return err
				}
			}
		case code.OpSetGlobal:
			{
				pos := int(binary.BigEndian.Uint16(ins[ip+1:]))
				vm.currentFrame().ip += 2
				if pos >= GlobalsSize {
					return fmt.Errorf("There can't be more than %d global variables", pos)
				}
				element := vm.pop()
				vm.globals[pos] = element
			}
		case code.OpJump:
			{
				pos := int(binary.BigEndian.Uint16(ins[ip+1:]))
				ip = pos - 1
			}
		case code.OpJumpNotTruthy:
			{
				pos := int(binary.BigEndian.Uint16(ins[ip+1:]))
				// Skip the 2 bytes of this operand
				vm.currentFrame().ip += 2
				condition := vm.pop()
				if !isTruthy(condition) {
					ip = pos - 1
				}
			}
		case code.OpGetLocal:
			{
				localIndex := byte(ins[ip+1])
				vm.currentFrame().ip++
				frame := vm.currentFrame()
				if err := vm.push(vm.stack[frame.basePointer+int(localIndex)]); err != nil {
					return nil
				}
			}
		case code.OpSetLocal:
			{
				localIndex := byte(ins[ip+1])
				vm.currentFrame().ip++
				frame := vm.currentFrame()
				vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
			}
		case code.OpCall:
			{
				fn, ok := vm.stack[vm.sp-1].(*object.CompiledFunction)
				if !ok {
					return fmt.Errorf("can't call expression, expected a function, got=%s", vm.stack[vm.sp-1].Type())
				}
				// Set the starting point for the function stack [..., fn, vm.sp..vm.sp+fn.NumLocals, stackOfTheFunction]
				frame := NewFrame(fn, vm.sp)
				vm.pushFrame(frame)
				vm.sp = frame.basePointer + fn.NumLocals
			}
		case code.OpNull:
			{
				if err := vm.push(Null); err != nil {
					return err
				}
			}
		case code.OpConstant:
			{
				idx := binary.BigEndian.Uint16(ins[ip+1:])
				vm.currentFrame().ip += 2
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
				case *object.String:
					{
						leftStr, ok := left.(*object.String)
						if !ok {
							return fmt.Errorf("expected %s, got: %s", right.Type(), left.Type())
						}
						if err := vm.push(&object.String{Value: leftStr.Value + rightObject.Value}); err != nil {
							return err
						}
					}
				}
			}
		case code.OpBang:
			{
				if err := vm.bangOperator(); err != nil {
					return err
				}
			}
		case code.OpMinus:
			{
				integer, err := vm.popIntegerObject()
				if err != nil {
					return err
				}
				integer.Value = -integer.Value
				vm.push(integer)
			}
		}
	}
	return nil
}

func (vm *VM) bangOperator() error {
	operand := vm.pop()
	switch operand {
	case True:
		{
			return vm.push(False)
		}
	case False, Null:
		{
			return vm.push(True)
		}
	}
	return vm.push(False)
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

// Frame

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}
