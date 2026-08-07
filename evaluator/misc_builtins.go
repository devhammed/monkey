package evaluator

import (
	"io/ioutil"
	"monkey/object"
	"monkey/typing"
)

func init() {
	builtins["require"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if err := typing.Check(
				"require",
				args,
				typing.ExactArgs(1),
				typing.WithTypes(object.STRING_OBJ),
			); err != nil {
				return newError(err.Error())
			}

			file := args[0].Inspect()
			data, err := ioutil.ReadFile(file)

			if err != nil {
				return newError("failed to require file: %s", err.Error())
			}

			env := object.NewEnvironment()

			evaluated := Run(string(data), file, FALSE, env)

			if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
				return newError(
					"error in required file (%s):\n %s",
					file,
					evaluated.Inspect(),
				)
			}

			return env.ExportedHash()
		},
	}

	builtins["len"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if err := typing.Check(
				"len",
				args,
				typing.ExactArgs(1),
			); err != nil {
				return newError(err.Error())
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	}

	builtins["type"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if err := typing.Check(
				"type",
				args,
				typing.ExactArgs(1),
			); err != nil {
				return newError(err.Error())
			}

			return &object.String{Value: string(args[0].Type())}
		},
	}
}
