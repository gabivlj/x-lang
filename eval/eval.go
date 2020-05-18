package eval

import (
	"xlang/ast"
	"xlang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func booleanToObject(boolean bool) *object.Boolean {
	if boolean {
		return TRUE
	}
	return FALSE
}

// Evaluator evaluates an AST of Xlang
type Evaluator struct {
	env *object.Environment
}

// NewEval returns a new evaluator of AST
func NewEval() *Evaluator {
	return &Evaluator{env: object.NewEnvironment()}
}

// ExtendEval returns a new evaluator with the environment you pass
func ExtendEval(env *object.Environment) *Evaluator {
	return &Evaluator{env: env}
}

// Eval evals an ast node.
func (e *Evaluator) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IndexExpression:
		{
			left := e.Eval(node.Left)
			if object.IsError(left) {
				return left
			}
			right := e.Eval(node.Right)
			if object.IsError(right) {
				return right
			}
			return e.evaluateIndex(left, right)
		}
	case *ast.ArrayLiteral:
		{
			elements := e.evalExpressions(node.Elements)
			if len(elements) == 1 && object.IsError(elements[0]) {
				return elements[0]
			}
			return &object.Array{Elements: elements}
		}

	case *ast.StringLiteral:
		{
			return &object.String{Value: node.Value}
		}

	case *ast.CallExpression:
		{
			function := e.Eval(node.Function)
			if object.IsError(function) {
				return function
			}
			parameters := e.evalExpressions(node.Arguments)
			if len(parameters) == 1 && object.IsError(parameters[0]) {
				return parameters[0]
			}
			return e.applyFunction(function, parameters)
		}
	case *ast.FunctionLiteral:
		{
			f := &object.Function{Parameters: node.Parameters, Body: node.Body, Env: e.env}
			return f
		}

	case *ast.LetStatement:
		{
			val := e.Eval(node.Value)
			if object.IsError(val) {
				return val
			}
			e.env.Set(node.Name.Value, val)

		}
	case *ast.Identifier:
		{
			return e.evalIdentifier(node)
		}
	case *ast.ReturnStatement:
		{
			val := e.Eval(node.ReturnValue)
			if object.IsError(val) {
				return val
			}
			return &object.ReturnValue{Value: val}
		}
	case *ast.BlockStatement:
		{
			return e.evalBlockStatement(node)
		}
	case *ast.IfExpression:
		{
			return e.evalIf(node)
		}
	case *ast.PrefixExpression:
		{
			value := e.Eval(node.Right)
			if object.IsError(value) {
				return value
			}
			return e.evalPrefixExpression(node.Operator, value)
		}
	case *ast.Program:
		{
			return e.evalProgramStatements(node.Statements)
		}

	case *ast.InfixExpression:
		{
			left := e.Eval(node.Left)
			if object.IsError(left) {
				return left
			}
			right := e.Eval(node.Right)
			if object.IsError(right) {
				return right
			}
			return e.evalInfixExpression(left, right, node.Operator)
		}

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.Boolean:
		{
			return booleanToObject(node.Value)
		}
	case *ast.IntegerLiteral:
		{
			return &object.Integer{Value: node.Value}
		}
	}
	return nil
}

func (e *Evaluator) applyFunction(fn object.Object, params []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		builtin, ok := fn.(*object.Builtin)
		if !ok {
			return object.NewError("Expected function, got %s instead", fn.Type())
		}
		return builtin.Fn(params...)
	}

	if len(params) != len(function.Parameters) {
		return object.NewError("Error, expected %d parameters, got %d", len(function.Parameters), len(params))
	}

	extendedEnvironment := e.newEnvironmentForFunction(function, params)
	eval := ExtendEval(extendedEnvironment)
	returnValue := eval.Eval(function.Body)
	tryUnwrapReturnValue, ok := returnValue.(*object.ReturnValue)
	if !ok {
		return returnValue
	}
	return tryUnwrapReturnValue
}

func (e *Evaluator) evaluateIndex(left object.Object, right object.Object) object.Object {
	switch obj := left.(type) {
	case *object.Array:
		{
			return e.evaluateArrayIndex(obj, right)
		}
	}
	return object.NewError("Unsupported index operation on type: %s", left.Type())
}

func (e *Evaluator) evaluateArrayIndex(array *object.Array, right object.Object) object.Object {
	idx, ok := right.(*object.Integer)
	if !ok {
		return object.NewError("Unsupported index on array of type: ", right.Type())
	}
	if int(idx.Value) >= len(array.Elements) || int(idx.Value) < 0 {
		return object.NewError("Array out of bounds, array size: %d, passed index: %d", len(array.Elements), idx.Value)
	}
	return array.Elements[idx.Value]
}

func (e *Evaluator) newEnvironmentForFunction(fn *object.Function, params []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for idx, param := range fn.Parameters {
		env.Set(param.Value, params[idx])
	}

	return env
}

func (e *Evaluator) evalExpressions(exps []ast.Expression) []object.Object {
	var result []object.Object
	for _, exp := range exps {
		evaluated := e.Eval(exp)
		if object.IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier) object.Object {
	val, ok := e.env.Get(node.Value)
	if !ok {
		if builtin, ok := builtins[node.Value]; ok {
			return builtin
		}

		// TODO: Maybe tell the user the most close variable?
		return object.NewError("Unknown variable %s", node.Value)
	}
	return val
}

func (e *Evaluator) evalIf(ifStatement *ast.IfExpression) object.Object {
	condition := e.Eval(ifStatement.Condition)
	if object.IsError(condition) {
		return condition
	}
	if NULL == condition || FALSE == condition {
		if ifStatement.Alternative == nil {
			return NULL
		}
		return e.Eval(ifStatement.Alternative)
	}

	return e.Eval(ifStatement.Consequence)
}

func (e *Evaluator) evalProgramStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = e.Eval(statement)
		if resultValue, ok := result.(*object.ReturnValue); ok {
			return resultValue.Value
		}
		if errorValue, isErr := result.(*object.Error); isErr {
			return errorValue
		}
	}
	return result
}

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = e.Eval(statement)

		if result != nil && (result.Type() == object.ReturnObject || result.Type() == object.ErrorObject) {
			return result
		}
	}
	return result
}

func (e *Evaluator) evalStringExpression(left *object.String, right *object.String, operator string) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	case "==":
		return booleanToObject(left.Value == right.Value)
	}
	return object.NewError("Error, unknown operator for strings '%s'", operator)
}

func (e *Evaluator) evalInfixExpression(left object.Object, right object.Object, operator string) object.Object {

	switch {
	case left.Type() != right.Type():
		{
			return object.NewError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
		}
	case left.Type() == object.StringObject:
		{
			return e.evalStringExpression(left.(*object.String), right.(*object.String), operator)
		}
	case left.Type() == object.IntegerObject:
		{
			return e.evalIntegerExpression(left.(*object.Integer), right.(*object.Integer), operator)
		}
	case operator == "==":
		{
			return booleanToObject(left == right)
		}
	case operator == "!=":
		{
			return booleanToObject(left != right)
		}
	}
	return object.NewError("Unknown operator: %s", operator)
}

func (e *Evaluator) evalIntegerExpression(left *object.Integer, right *object.Integer, operator string) object.Object {
	switch operator {
	case "*":
		{
			return &object.Integer{Value: left.Value * right.Value}
		}
	case "/":
		{
			return &object.Integer{Value: left.Value / right.Value}
		}
	case "+":
		{
			return &object.Integer{Value: left.Value + right.Value}
		}
	case "-":
		{
			{
				return &object.Integer{Value: left.Value - right.Value}
			}
		}
	case "%":
		{
			return &object.Integer{Value: left.Value % right.Value}
		}

	case ">":
		{
			return booleanToObject(left.Value > right.Value)
		}
	case "<":
		{
			return booleanToObject(left.Value < right.Value)
		}
	case "<=":
		{
			return booleanToObject(left.Value <= right.Value)
		}
	case ">=":
		{
			return booleanToObject(left.Value >= right.Value)
		}
	case "==":
		{
			return booleanToObject(left.Value == right.Value)
		}
	case "!=":
		{
			return booleanToObject(left.Value != right.Value)
		}
	}
	// todo: Throw err
	return object.NewError("Unknown operator: %s", operator)
}

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return e.evalBangOperatorRight(right)
	case "-":
		return e.evalMinusOperatorRight(right)
	}
	return object.NewError("Unknown prefix operator: %s", operator)
}

func (e *Evaluator) evalMinusOperatorRight(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		{
			right.Value = -right.Value
			return right
		}
	default:
		{
			return object.NewError("Mismatch type left operator -%s", right.Inspect())
		}
	}
}

func (e *Evaluator) evalBangOperatorRight(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
