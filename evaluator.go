package monkey_interpreter

var (
	FALSE_OBJ = &BooleanObject{Value: false}
	TRUE_OBJ  = &BooleanObject{Value: true}
)

func Eval(node Node) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalStatements(node.Statements)
	case *ExpressionStatement:
		return Eval(node.Expression)
		// Expressions
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	case *BooleanLiteral:
		if node.Value {
			return TRUE_OBJ
		}
		return FALSE_OBJ
	}
	return nil
}

func evalStatements(stmts []Statement) Object {
	var result Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}
