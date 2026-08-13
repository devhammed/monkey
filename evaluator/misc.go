package evaluator

import (
	"monkey/object"
	"monkey/typing"
	"os"
	"path/filepath"
)

var (
	NULL    = &object.Null{}
	TRUE    = &object.Boolean{Value: true}
	FALSE   = &object.Boolean{Value: false}
	VERSION = &object.String{Value: "v0.2.7"}
)

func init() {

	superGlobals["VERSION"] = VERSION

	superGlobals["ARGV"] = ARGV

	builtins["require"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"require",
				args,
				typing.ExactArgs(1),
				typing.WithTypes(object.STRING_OBJ),
			); err != nil {
				return newError("%s", err.Error())
			}

			file := args[0].Inspect()
			abs, err := filepath.Abs(file)
			if err != nil {
				return newError("failed to get absolute path for file: %s", err.Error())
			}
			data, err := os.ReadFile(abs)

			if err != nil {
				return newError("failed to require file: %s", err.Error())
			}

			moduleEnv := object.NewModuleEnvironment(env)

			evaluated := Run(string(data), abs, filepath.Dir(abs), false, moduleEnv)

			if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
				return newError(
					"error in required file (%s):\n %s",
					file,
					evaluated.Inspect(),
				)
			}

			return moduleEnv.ExportedHash()
		},
	}

	builtins["len"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"len",
				args,
				typing.ExactArgs(1),
			); err != nil {
				return newError("%s", err.Error())
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
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"type",
				args,
				typing.ExactArgs(1),
			); err != nil {
				return newError("%s", err.Error())
			}

			return &object.String{Value: string(args[0].Type())}
		},
	}
}
