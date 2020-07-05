package compiler

import (
	"fmt"
	"xlang/ast"
	"xlang/code"
	"xlang/object"
)

// EmittedInstruction Keeps track of an emitted instruction
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// Compiler contains the instructions and constants
type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// New returns a new compiler
func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i, insByte := range newInstruction {
		c.instructions[i+pos] = insByte
	}
}

func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

// Compile saves in the compiler the instructions that the ast node produces
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.IfExpression:
		{
			err := c.Compile(node.Condition)
			if err != nil {
				return err
			}
			pos := c.emit(code.OpJumpNotTruthy, 9999)
			err = c.Compile(node.Consequence)
			if err != nil {
				return err
			}
			if c.lastInstruction.Opcode == code.OpPop {
				c.instructions = c.instructions[:c.lastInstruction.Position]
				c.lastInstruction = c.previousInstruction
			}
			posOfJump := c.emit(code.OpJump, 9999)
			c.changeOperand(pos, len(c.instructions))
			if node.Alternative != nil {
				c.Compile(node.Alternative)
				if c.lastInstruction.Opcode == code.OpPop {
					c.instructions = c.instructions[:c.lastInstruction.Position]
					c.lastInstruction = c.previousInstruction
				}
				c.changeOperand(posOfJump, len(c.instructions))
				return nil
			}
			// If there is no alternative, "fake" it
			c.emit(code.OpNull)
			c.changeOperand(posOfJump, len(c.instructions))
		}
	case *ast.BlockStatement:
		{
			for _, s := range node.Statements {
				if err := c.Compile(s); err != nil {
					return err
				}
			}
		}
	case *ast.Program:
		{
			for _, s := range node.Statements {
				err := c.Compile(s)
				if err != nil {
					return err
				}
			}
		}
	case *ast.PrefixExpression:
		{
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			switch node.Operator {
			case "!":
				c.emit(code.OpBang)
			case "-":
				c.emit(code.OpMinus)
			default:
				return fmt.Errorf("unknown prefix operator: %s", node.Operator)
			}
		}
	case *ast.ExpressionStatement:
		{
			err := c.Compile(node.Expression)
			if err != nil {
				return err
			}
			c.emit(code.OpPop)
		}
	case *ast.Boolean:
		{
			if node.Value {
				c.emit(code.OpTrue)
			} else {
				c.emit(code.OpFalse)
			}
		}
	case *ast.InfixExpression:
		{
			// Change the order of operations in case
			nodeToUseForLeft := node.Left
			nodeToUseForRight := node.Right
			if node.Operator == "<" {
				nodeToUseForLeft = node.Right
				nodeToUseForRight = node.Left
			}
			err := c.Compile(nodeToUseForLeft)
			if err != nil {
				return err
			}
			err = c.Compile(nodeToUseForRight)
			if err != nil {
				return err
			}
			switch node.Operator {
			case "+":
				c.emit(code.OpAdd)
			case "-":
				c.emit(code.OpSub)
			case "*":
				c.emit(code.OpMul)
			case "/":
				c.emit(code.OpDiv)
			case "<":
				// Check beginning of case, we change the order of operators
				c.emit(code.OpGreaterThan)
			case ">":
				c.emit(code.OpGreaterThan)
			case "==":
				c.emit(code.OpEqual)
			case "!=":
				c.emit(code.OpNotEqual)

			default:
				return fmt.Errorf("unknown operator %s", node.Operator)
			}
		}
	case *ast.IntegerLiteral:
		{
			// Create integer objet
			integer := &object.Integer{Value: node.Value}
			// Create the VM code for adding a constant into the stack
			c.emit(code.OpConstant, c.addConstant(integer))
		}
	}
	return nil
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewIns := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewIns
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// Bytecode contains all the bytecode and constants
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// Bytecode returns the bytecode of the compiler
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
