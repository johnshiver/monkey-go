package monkey_interpreter

import (
	"fmt"
)

var (
	NULL_OBJ  = &Null{}
	FALSE_OBJ = &BooleanObject{Value: false}
	TRUE_OBJ  = &BooleanObject{Value: true}
)

func Eval(node Node, env *Environment) Object {
	switch currNode := node.(type) {
	// Statements
	case *Program:
		return evalProgram(currNode, env)
	case *ExpressionStatement:
		return Eval(currNode.Expression, env)
	case *BlockStatement:
		return evalBlockStatement(currNode, env)
	case *IfExpression:
		return evalIfExpression(currNode, env)
	case *ReturnStatement:
		val := Eval(currNode.ReturnValue, env)
		return &ReturnValue{Value: val}
	case *LetStatement:
		// evaluate let expression
		val := Eval(currNode.Value, env)

		// if there is an error evaluating expression, return the error
		if isError(val) {
			return val
		}
		env.Set(currNode.Name.Value, val)

	case *Identifier:
		return evalIdentifier(currNode, env)

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: currNode.Value}
	case *BooleanLiteral:
		return nativeBoolToBooleanObject(currNode.Value)
	case *PrefixExpression:
		right := Eval(currNode.Right, env)
		return evalPrefixExpression(currNode.Operator, right)
	case *InfixExpression:
		left := Eval(currNode.Left, env)
		right := Eval(currNode.Right, env)
		return evalInfixExpression(currNode.Operator, left, right)
	case *FunctionLiteral:
		params := currNode.Parameters
		body := currNode.Body
		return &Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
	case *CallExpression:
		function := Eval(currNode.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(currNode.Arguments, env)
		// if there is an error on the previous call, this is how to check for it
		// TODO: im not a huge fan of this design, lets revisi
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

func evalProgram(program *Program, env *Environment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch resultType := result.(type) {
		case *ReturnValue:
			return resultType.Value
		case *Error:
			return resultType
		}
	}
	return result
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE_OBJ:
		return FALSE_OBJ
	case FALSE_OBJ:
		return TRUE_OBJ
	case NULL_OBJ:
		return TRUE_OBJ
	default:
		return FALSE_OBJ
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INT_OBJ_TYPE {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INT_OBJ_TYPE && right.Type() == INT_OBJ_TYPE:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value
	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *BooleanObject {
	if input {
		return TRUE_OBJ
	}
	return FALSE_OBJ
}

func evalIfExpression(ie *IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL_OBJ
	}
}
func isTruthy(obj Object) bool {
	switch obj {
	case NULL_OBJ:
		return false
	case TRUE_OBJ:
		return true
	case FALSE_OBJ:
		return false
	default:
		return true
	}
}

func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ_TYPE || rt == ERROR_OBJ_TYPE {
				return result
			}
		}
	}
	return result
}

func evalIdentifier(node *Identifier, env *Environment) Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalExpressions(
	exps []Expression, env *Environment,
) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn Object, args []Object) Object {
	function, ok := fn.(*Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
	return obj.Type() == ERROR_OBJ_TYPE
}
