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
	case *StringLiteral:
		return &String{Value: currNode.Value}
	case *BooleanLiteral:
		return nativeBoolToBooleanObject(currNode.Value)
	case *ArrayLiteral:
		elements := evalExpressions(currNode.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}
	case *HashLiteral:
		return evalHashLiteral(currNode, env)
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
	case *IndexExpression:
		left := Eval(currNode.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(currNode.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
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
	case left.Type() == STRING_OBJ_TYPE && right.Type() == STRING_OBJ_TYPE:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}
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

func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ_TYPE && index.Type() == INT_OBJ_TYPE:
		return evalArrayIndexExpression(left, index)
	case left.Type() == HASH_OBJ_TYPE:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL_OBJ
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index Object) Object {
	hashObject := hash.(*Hash)
	key, ok := index.(Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL_OBJ
	}
	return pair.Value
}

func evalHashLiteral(
	node *HashLiteral, env *Environment,
) Object {

	pairs := make(map[HashKey]HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key: key, Value: value}
	}
	return &Hash{Pairs: pairs}
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
	if ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
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
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
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
