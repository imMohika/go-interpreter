package eval

import "go-interpreter/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fun: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"head": {
		Fun: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) == 0 {
					return NULL
				}
				return &object.String{Value: string(arg.Value[0])}
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument to `head` not supported, got %s", args[0].Type())
			}
		},
	},
	"tail": {
		Fun: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) == 0 {
					return NULL
				}
				return &object.String{Value: string(arg.Value[1:len(arg.Value)])}
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return NULL
				}
				tailElements := make([]object.Object, length-1)
				copy(tailElements, arg.Elements[1:length])
				return &object.Array{Elements: tailElements}
			default:
				return newError("argument to `tail` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fun: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.String:
				length := len(arg.Value)
				if length == 0 {
					return NULL
				}
				return &object.String{Value: string(arg.Value[length-1])}
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return NULL
				}
				return arg.Elements[length-1]
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fun: func(args ...object.Object) object.Object {
			if err := checkArgsLen(2, args...); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return NULL
				}
				newElements := make([]object.Object, length+1)
				copy(newElements, arg.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `push` must be Array, got %s", args[0].Type())
			}
		},
	},
}

func checkArgsLen(argsLen int, args ...object.Object) *object.Error {
	if len(args) != argsLen {
		return newError("wrong number of arguments, got=%d, want=%d", len(args), argsLen)
	}
	return nil
}
