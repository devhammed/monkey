package evaluator

import (
	"monkey/object"
	"monkey/typing"
)

func init() {
	builtins["range"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"range",
				args,
				typing.RangeOfArgs(2, 3),
				typing.AllOfType(object.INTEGER_OBJ),
			); err != nil {
				return newError(err.Error())
			}

			step := int64(1)
			end := args[1].(*object.Integer)
			start := args[0].(*object.Integer)

			if len(args) == 3 {
				s := args[2].(*object.Integer)
				step = s.Value
			}

			i := start.Value
			var arr []object.Object

			for i < end.Value {
				arr = append(arr, &object.Integer{Value: i})
				i = i + step
			}

			return &object.Array{Elements: arr}
		},
	}

	builtins["array_first"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_first", args, typing.ExactArgs(1), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	}

	builtins["array_last"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_last", args, typing.ExactArgs(1), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}
			return NULL
		},
	}

	builtins["array_rest"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_rest", args, typing.ExactArgs(1), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				elements := make([]object.Object, len(arr.Elements)-1)
				copy(elements, arr.Elements[1:])
				return &object.Array{Elements: elements}
			}
			return NULL
		},
	}

	builtins["array_push"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_push", args, typing.ExactArgs(2), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}

			arr := args[0].(*object.Array)
			elements := make([]object.Object, len(arr.Elements)+1)
			copy(elements, arr.Elements)
			elements[len(arr.Elements)] = args[1]
			arr.Elements = elements
			return NULL
		},
	}

	builtins["array_map"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_map", args, typing.ExactArgs(2), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}
			if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
				return newError("second argument to `array_map` must be FUNCTION, got %s", args[1].Type())
			}

			arr := args[0].(*object.Array)
			elements := make([]object.Object, len(arr.Elements))
			for i, element := range arr.Elements {
				elements[i] = applyFunction(args[1], []object.Object{element, &object.Integer{Value: int64(i)}}, env)
			}
			return &object.Array{Elements: elements}
		},
	}

	builtins["array_each"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_each", args, typing.ExactArgs(2), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}
			if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
				return newError("second argument to `array_each` must be FUNCTION, got %s", args[1].Type())
			}

			for i, element := range args[0].(*object.Array).Elements {
				applyFunction(args[1], []object.Object{element, &object.Integer{Value: int64(i)}}, env)
			}
			return NULL
		},
	}

	builtins["array_reduce"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"array_reduce",
				args,
				typing.ExactArgs(3),
				typing.WithTypes(object.ARRAY_OBJ),
			); err != nil {
				return newError(err.Error())
			}

			if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
				return newError("second argument to `array_reduce` must be FUNCTION, got %s",
					args[1].Type())
			}

			arr := args[0].(*object.Array)
			elements := arr.Elements
			acc := args[2]

			for i := 0; i < len(elements); i++ {
				fnArgs := []object.Object{acc, elements[i], &object.Integer{Value: int64(i)}}

				acc = applyFunction(acc, fnArgs, env)
			}

			return acc
		},
	}

	builtins["array_copy"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("array_copy", args, typing.ExactArgs(1), typing.WithTypes(object.ARRAY_OBJ)); err != nil {
				return newError(err.Error())
			}

			elements := make([]object.Object, len(args[0].(*object.Array).Elements))
			copy(elements, args[0].(*object.Array).Elements)
			return &object.Array{Elements: elements}
		},
	}
}
