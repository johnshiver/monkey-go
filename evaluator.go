package monkey_interpreter

func Eval(node Node) Object {
	switch node := node.(type) {
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	}
	return nil
}
