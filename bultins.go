package monkey_interpreter

import "fmt"

var builtins = map[string]*Builtin{
	"len": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ_TYPE {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*Array)

			// NOTE: i would prefer to reuse the original array
			// the book makes a new copy

			//length := len(arr.Elements)
			//newElements := make([]Object, length+1, length+1)
			//copy(newElements, arr.Elements)
			//newElements[length] = args[1]
			arr.Elements = append(arr.Elements, args[1])
			//return &Array{Elements: newElements}
			return arr
		},
	},
	"puts": {
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL_OBJ
		},
	},
}
