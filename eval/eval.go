package eval

import (
	"fmt"
	"go-interpreter/ast"
	"go-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Value, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return booleanFromNativeBool(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	default:
		return newError("invalid node: %s", node.String())
	}
	return nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}

	return false
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		if isError(result) {
			return result
		}

		if result.Type() == object.RETURN_VALUE_OBJECT || result.Type() == object.ERROR_OBJECT {
			return result
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return value
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	default:
		return TRUE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return newError("invalid usage of `-` operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	switch {
	case left.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJECT:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJECT:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}

	case "<":
		return booleanFromNativeBool(leftValue < rightValue)
	case "<=":
		return booleanFromNativeBool(leftValue <= rightValue)
	case ">":
		return booleanFromNativeBool(leftValue > rightValue)
	case ">=":
		return booleanFromNativeBool(leftValue >= rightValue)
	case "==":
		return booleanFromNativeBool(leftValue == rightValue)
	case "!=":
		return booleanFromNativeBool(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return booleanFromNativeBool(leftValue == rightValue)
	case "!=":
		return booleanFromNativeBool(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return booleanFromNativeBool(leftValue == rightValue)
	case "!=":
		return booleanFromNativeBool(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition == TRUE {
		return Eval(node.Consequence, env)
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	// todo))
	return NULL
}

func booleanFromNativeBool(value bool) object.Object {
	switch value {
	case true:
		return TRUE
	default:
		return FALSE
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalExpressions(arguments []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, 0, len(arguments))
	for _, expr := range arguments {
		evaluated := Eval(expr, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(function object.Object, args []object.Object) object.Object {
	fun, ok := function.(*object.Function)
	if !ok {
		return newError("not a function: %s", function.Type())
	}

	innerEnv := extendFunctionEnv(fun, args)
	evaluated := Eval(fun.Body, innerEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fun *object.Function, args []object.Object) *object.Environment {
	env := object.NewInnerEnvironment(fun.Env)
	for idx, param := range fun.Parameters {
		env.Set(param.Value, args[idx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
