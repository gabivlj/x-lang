package compiler

import (
	"fmt"
	"xlang/ast"
	"xlang/code"
	"xlang/object"
)

// Compiler contains the instructions and constants
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// New returns a new compiler
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile saves in the compiler the instructions that the ast node produces
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		{
			for _, s := range node.Statements {
				err := c.Compile(s)
				if err != nil {
					return err
				}
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
