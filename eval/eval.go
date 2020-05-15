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
	case *ast.PrefixExpression:
		{
			value := Eval(node.Right)
			return evalPrefixExpression(node.Operator, value)
		}
	case *ast.Program:
		{
			return evalStatements(node.Statements)
		}

	case *ast.InfixExpression:
		{
			left := Eval(node.Left)
			right := Eval(node.Right)
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

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}

func evalInfixExpression(left object.Object, right object.Object, operator string) object.Object {
	if left.Type() != right.Type() {
		return NULL
	}
	if left.Type() == object.IntegerObject {
		return evalIntegerExpression(left.(*object.Integer), right.(*object.Integer), operator)
	}
	return NULL
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
	}
	// todo: Throw err
	return NULL
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorRight(right)
	case "-":
		return evalMinusOperatorRight(right)
	}
	return NULL
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
			return NULL
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
