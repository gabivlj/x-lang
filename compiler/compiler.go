package compiler

import (
	"fmt"
	"sort"
	"xlang/ast"
	"xlang/code"
	"xlang/object"
)

// EmittedInstruction Keeps track of an emitted instruction
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// CompilationScope contains everything about the current scope
type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	name                string
}

// Compiler contains the instructions and constants
type Compiler struct {
	constants []object.Object

	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIndex  int
}

// New returns a new compiler
func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	table := NewSymbolTable()
	for i, fn := range object.GetBuiltins() {
		table.DefineBuiltin(i, fn.Name)
	}
	return &Compiler{
		constants:   []object.Object{},
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		symbolTable: table,
	}
}

// NewWithState returns a new compiler with a symbol table or constants inserted already
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (c *Compiler) enterScope(name string) {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		name:                name,
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return instructions
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) currentScope() *CompilationScope {
	return &c.scopes[c.scopeIndex]
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()
	for i, insByte := range newInstruction {
		ins[i+pos] = insByte
	}
}

func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// Compile saves in the compiler the instructions that the ast node produces
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.IndexExpression:
		{
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			c.emit(code.OpIndex)
		}
	case *ast.HashLiteral:
		{
			keys := make([]ast.Expression, 0, len(node.Pairs))
			for k := range node.Pairs {
				keys = append(keys, k)
			}
			// Make a consistent order in keys
			sort.Slice(keys, func(i, j int) bool {
				return keys[i].String() < keys[j].String()
			})
			for _, key := range keys {
				value := node.Pairs[key]
				err := c.Compile(key)
				if err != nil {
					return err
				}
				errVal := c.Compile(value)
				if errVal != nil {
					return errVal
				}
			}
			c.emit(code.OpHash, len(node.Pairs)*2)
		}
	case *ast.ArrayLiteral:
		{
			for _, exp := range node.Elements {
				if err := c.Compile(exp); err != nil {
					return err
				}
			}
			c.emit(code.OpArray, len(node.Elements))
		}
	case *ast.CallExpression:
		{
			if err := c.Compile(node.Function); err != nil {
				return err
			}
			for _, a := range node.Arguments {
				if err := c.Compile(a); err != nil {
					return err
				}
			}
			c.emit(code.OpCall, len(node.Arguments))
		}
	case *ast.StringLiteral:
		{
			str := node.Value
			pos := c.addConstant(&object.String{Value: str})
			c.emit(code.OpConstant, pos)
		}
	case *ast.Identifier:
		{
			symbol, ok := c.symbolTable.Resolve(node.Value)
			if node.Value == c.currentScope().name && symbol.Scope == FreeScope {
				c.emit(code.OpCurrentClosure)
				return nil
			}
			if !ok {
				return fmt.Errorf("undefined variable=%s que", node.Value)
			}
			c.emit(c.getCodeScope(&symbol), symbol.Index)
		}
	case *ast.LetStatement:
		{
			symbol := c.symbolTable.Define(node.Name.Value)
			if err := c.Compile(node.Value); err != nil {
				return err
			}

			c.emit(c.setCodeScope(&symbol), symbol.Index)
		}
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
			scope := c.currentScope()
			if scope.lastInstruction.Opcode == code.OpPop {
				last := c.scopes[c.scopeIndex].lastInstruction
				previous := c.scopes[c.scopeIndex].previousInstruction

				old := c.currentInstructions()
				new := old[:last.Position]

				c.scopes[c.scopeIndex].instructions = new
				c.scopes[c.scopeIndex].lastInstruction = previous
			}
			posOfJump := c.emit(code.OpJump, 9999)
			c.changeOperand(pos, len(scope.instructions))
			if node.Alternative != nil {
				c.Compile(node.Alternative)
				if scope.lastInstruction.Opcode == code.OpPop {
					last := c.scopes[c.scopeIndex].lastInstruction
					previous := c.scopes[c.scopeIndex].previousInstruction

					old := c.currentInstructions()
					new := old[:last.Position]

					c.scopes[c.scopeIndex].instructions = new
					c.scopes[c.scopeIndex].lastInstruction = previous
				}
				c.changeOperand(posOfJump, len(scope.instructions))
				return nil
			}
			// If there is no alternative, "fake" it
			c.emit(code.OpNull)
			c.changeOperand(posOfJump, len(scope.instructions))
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

	case *ast.ReturnStatement:
		{
			if err := c.Compile(node.ReturnValue); err != nil {
				return err
			}

			c.emit(code.OpReturnValue)

		}

	case *ast.FunctionLiteral:
		{
			c.enterScope(node.Name)
			// All the parameters will be already in the stack so we won't need to worry about
			// it. [..., FUNCTION_LITERAL, ...argumentsOnStack, ...localVariables]
			// OpCodes: [OpConstants(FUNCTION), ...OpConstants | OpGet..., OpCall]
			for _, p := range node.Parameters {
				c.symbolTable.Define(p.Value)
			}
			if err := c.Compile(node.Body); err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.replaceInstruction(c.currentScope().lastInstruction.Position, code.Make(code.OpReturnValue))
				c.currentScope().lastInstruction.Opcode = code.OpReturnValue
			}
			if !c.lastInstructionIs(code.OpReturnValue) {
				c.emit(code.OpReturn)
			}
			numLocals := c.symbolTable.numDefinitions
			freeSymbols := c.symbolTable.FreeSymbols
			ins := c.leaveScope()
			for _, s := range freeSymbols {
				c.emit(c.getCodeScope(&s), s.Index)
			}
			compiledFn := &object.CompiledFunction{Instructions: ins, NumLocals: numLocals, NumParameters: len(node.Parameters)}
			c.emit(code.OpClosure, c.addConstant(compiledFn), len(freeSymbols))
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
	posNewIns := len(c.currentInstructions())
	c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), ins...)
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
	Table        *SymbolTable
}

// Bytecode returns the bytecode of the compiler
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
		Table:        c.symbolTable,
	}
}

func (c *Compiler) setCodeScope(symbol *Symbol) code.Opcode {
	codeToUse := code.OpSetGlobal
	if symbol.Scope == LocalScope {
		codeToUse = code.OpSetLocal
	}
	return codeToUse
}

func (c *Compiler) getCodeScope(symbol *Symbol) code.Opcode {
	codeToUse := code.OpGetGlobal
	if symbol.Scope == LocalScope {
		codeToUse = code.OpGetLocal
	} else if symbol.Scope == BuiltinScope {
		codeToUse = code.OpGetBuiltin
	} else if symbol.Scope == FreeScope {
		codeToUse = code.OpGetFree
	}
	return codeToUse
}
