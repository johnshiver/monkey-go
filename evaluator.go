package monkey_interpreter

var (
	NULL_OBJ  = &Null{}
	FALSE_OBJ = &BooleanObject{Value: false}
	TRUE_OBJ  = &BooleanObject{Value: true}
)

func Eval(node Node) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalProgram(node)
	case *ExpressionStatement:
		return Eval(node.Expression)
	case *BlockStatement:
		return evalBlockStatement(node)
	case *IfExpression:
		return evalIfExpression(node)
	case *ReturnStatement:
		val := Eval(node.ReturnValue)
		return &ReturnValue{Value: val}

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	case *BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalProgram(program *Program) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		// if we're at return statement, short circuit return the value
		if returnVal, ok := result.(*ReturnValue); ok {
			return returnVal.Value
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
		return NULL_OBJ
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
		return NULL_OBJ
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
	default:
		return NULL_OBJ
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
		return NULL_OBJ
	}
}

func nativeBoolToBooleanObject(input bool) *BooleanObject {
	if input {
		return TRUE_OBJ
	}
	return FALSE_OBJ
}

func evalIfExpression(ie *IfExpression) Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
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

func evalBlockStatement(block *BlockStatement) Object {
	var result Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == RETURN_VALUE_OBJ_TYPE {
			return result
		}
	}
	return result
}
