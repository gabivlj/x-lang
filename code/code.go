package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions is the list of instructions that are passed to the VM
type Instructions []byte

// Opcode is a operation like PUSH
type Opcode byte

const (
	// OpConstant is a constant address
	OpConstant Opcode = iota
	// OpAdd is a operation for adding 2 numbers
	OpAdd
	// OpPop is an operation that tells the VM to pop the topmost element off the stack
	OpPop
	// OpSub operation for substract
	OpSub
	// OpMul operation for multiplication
	OpMul
	// OpDiv operation for dividing
	OpDiv
	// OpTrue pushes a true object into the stack
	OpTrue
	// OpFalse pushes a false object into the stack
	OpFalse
	// OpEqual makes a == comparison
	OpEqual
	// OpNotEqual makes a != comparison
	OpNotEqual
	// OpGreaterThan makes a > comparison
	OpGreaterThan
	// OpMinus -
	OpMinus
	// OpBang !
	OpBang
	// OpJumpNotTruthy jumps if the last element in the stack is not true
	OpJumpNotTruthy
	// OpJump jumps to the desired location
	OpJump
	// OpNull tells the VM to put Null in the stack
	OpNull
	// OpGetGlobal tells the VM to retrieve a global variable
	OpGetGlobal
	// OpSetGlobal tells the VM to set the last variable pushed into the stack as a global variable
	OpSetGlobal
	// OpArray tells the VM to add X elements from the stack into an array
	OpArray
	// OpHash tells the vm to add X key and values that are inside the stack into a hashmap
	OpHash
	// OpIndex tells the vm to add into the stack an element of the indexed
	OpIndex
	// OpCall tells the VM to call the function that is the last inserted element in the stack
	OpCall
	// OpReturnValue tells the VM to return from the function with a return value of the top of the stack
	OpReturnValue
	// OpReturn tells the VM to stop the execution of the current function returning Null
	OpReturn
	// OpGetLocal tells the VM to get a local variable
	OpGetLocal
	// OpSetLocal tells the VM to set a local variable
	OpSetLocal
	// OpGetBuiltin tells the VM to get a builtin variable/function in the position X
	OpGetBuiltin
	// OpClosure tells the specified X constant function to convert it into a closure and Y specifies how many free varibles are on the stack and need to be transfered into the closure.
	OpClosure
	// OpGetFree tells the VM to retrieve a free variable
	OpGetFree
	// OpCurrentClosure tells the VM to push to the stack the current closure
	OpCurrentClosure
)

// Definition is the definition of a operand
type Definition struct {
	Name          string
	OperandWidths []int
}

// definitions of operands
var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	// 1 being the number of arguments
	OpCall:        {"OpCall", []int{1}},
	OpReturnValue: {"OpReturnValue", []int{}},
	OpReturn:      {"OpReturn", []int{}},
	// The direction in the stack
	OpGetLocal: {"OpGetLocal", []int{1}},
	// Puts the current top element of the stack to the specified direction in the stack (initial:vm.sp + dir)
	OpSetLocal:   {"OpSetLocal", []int{1}},
	OpGetBuiltin: {"OpGetBuiltin", []int{1}},
	// OpClosure before it's emmited must have all the free variables loaded into the stack
	// with OpConstants/OpGet...
	OpClosure:        {"OpClosure", []int{2, 1}},
	OpGetFree:        {"OpGetFree", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

// Lookup an operand in the definition table
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Make returns (in big endian encoding) the byte slice of an operation with its operand, basically a instruction.
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	// offset starts with 1 because the operand!
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		// [...offset, byte, byte]
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d in %s %v\n",
			len(operands), operandCount, def.Name, operands)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// func (ins Instructions) String() string {
// 	if len(ins) == 0 {
// 		return ""
// 	}
// 	def, err := Lookup(ins[0])
// 	if err != nil {
// 		return ""
// 	}
// 	s := strings.Builder{}
// 	s.Grow(len(ins) * 10)
// 	offset := 0
// 	for i := 1; i < len(ins); i++ {
// 		if offset == 0 {
// 			s.WriteString(fmt.Sprintf("%04d ", i-1))
// 			s.WriteString(def.Name)
// 		}
// 		if len(def.OperandWidths) > 0 {
// 			str := opToString(i, def.OperandWidths[offset]+i, def.OperandWidths[offset], ins)
// 			s.WriteString(str)
// 			offset++
// 		}
// 		if offset >= len(def.OperandWidths) {
// 			// Go to next operand (Taking into mind that "i" is next to the current operator byte)
// 			if len(def.OperandWidths) > 0 {
// 				i = def.OperandWidths[offset-1] + i
// 			}
// 			if i >= len(ins) {
// 				break
// 			}
// 			offset = 0
// 			s.WriteByte('\n')
// 			def, _ = Lookup(ins[i])
// 			if def == nil {
// 				continue
// 			}
// 			if len(def.OperandWidths) == 0 && i == len(ins)-1 {
// 				s.WriteString(fmt.Sprintf("%04d ", i))
// 				s.WriteString(def.Name)
// 			}
// 		}
// 	}
// 	return strings.TrimSpace(s.String())
// }

func opToString(start, end, width int, ins Instructions) string {
	switch width {
	case 1:
		{
			return fmt.Sprintf(" %d", ins[start])
		}
	case 0:
		{
			return ""
		}
	case 2:
		{
			return fmt.Sprintf(" %d", binary.BigEndian.Uint16(ins[start:end]))
		}
	}
	return ""
}

// ReadOperands reverts the code of an operation
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, 0, len(def.OperandWidths))
	offset := 0
	for _, w := range def.OperandWidths {
		switch w {
		case 2:
			operands = append(operands, int(binary.BigEndian.Uint16(ins[offset:])))
		case 1:
			operands = append(operands, int(ins[offset]))
		}
		offset += w
	}
	return operands, offset
}
