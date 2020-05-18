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

// Eval evals the ast tree
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.ReturnStatement:
		{
			val := Eval(node.ReturnValue)
			if object.IsError(val) {
				return val
			}
			return &object.ReturnValue{Value: val}
		}
	case *ast.BlockStatement:
		{
			return evalBlockStatement(node)
		}
	case *ast.IfExpression:
		{
			return evalIf(node)
		}
	case *ast.PrefixExpression:
		{
			value := Eval(node.Right)
			if object.IsError(value) {
				return value
			}
			return evalPrefixExpression(node.Operator, value)
		}
	case *ast.Program:
		{
			return evalProgramStatements(node.Statements)
		}

	case *ast.InfixExpression:
		{
			left := Eval(node.Left)
			if object.IsError(left) {
				return left
			}
			right := Eval(node.Right)
			if object.IsError(right) {
				return right
			}
			return evalInfixExpression(left, right, node.Operator)
		}

	case *ast.ExpressionStatement:
		return Eval(node.Expression)
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

func evalIf(ifStatement *ast.IfExpression) object.Object {
	condition := Eval(ifStatement.Condition)
	if object.IsError(condition) {
		return condition
	}
	if NULL == condition || FALSE == condition {
		if ifStatement.Alternative == nil {
			return NULL
		}
		return Eval(ifStatement.Alternative)
	}

	return Eval(ifStatement.Consequence)
}

func evalProgramStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
		if resultValue, ok := result.(*object.ReturnValue); ok {
			return resultValue.Value
		}
		if errorValue, isErr := result.(*object.Error); isErr {
			return errorValue
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil && (result.Type() == object.ReturnObject || result.Type() == object.ErrorObject) {
			return result
		}
	}
	return result
}

func evalInfixExpression(left object.Object, right object.Object, operator string) object.Object {

	switch {
	case left.Type() != right.Type():
		{
			return object.NewError("Type mismatch: %s  %s", left.Type(), operator, right.Type())
		}
	case left.Type() == object.IntegerObject:
		{
			return evalIntegerExpression(left.(*object.Integer), right.(*object.Integer), operator)
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

func evalIntegerExpression(left *object.Integer, right *object.Integer, operator string) object.Object {
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorRight(right)
	case "-":
		return evalMinusOperatorRight(right)
	}
	return object.NewError("Unknown prefix operator: %s", operator)
}

func evalMinusOperatorRight(right object.Object) object.Object {
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

func evalBangOperatorRight(right object.Object) object.Object {
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
