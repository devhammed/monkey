package evaluator

import (
	"io"
	"monkey/object"
	"monkey/typing"
)

func init() {
	builtins["print"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"print",
				args,
				typing.MinimumArgs(1),
			); err != nil {
				return newError("%s", err.Error())
			}

			stdout, ok := env.Get("STDOUT")
			if !ok {
				return newError("print: STDOUT is not defined")
			}

			resource, ok := stdout.Value.(*object.Resource)
			if !ok {
				return newError("print: STDOUT is not a resource")
			}

			writer, ok := resource.Handle.(io.Writer)
			if !ok {
				return newError("print: STDOUT is not writable")
			}

			for _, arg := range args {
				if _, err := io.WriteString(writer, arg.Inspect()); err != nil {
					return newError("print: %s", err)
				}
			}

			return NULL
		},
	}

	builtins["println"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"println",
				args,
				typing.MinimumArgs(1),
			); err != nil {
				return newError("%s", err.Error())
			}

			stdout, ok := env.Get("STDOUT")
			if !ok {
				return newError("println: STDOUT is not defined")
			}

			resource, ok := stdout.Value.(*object.Resource)
			if !ok {
				return newError("println: STDOUT is not a resource")
			}

			writer, ok := resource.Handle.(io.Writer)
			if !ok {
				return newError("println: STDOUT is not writable")
			}

			for _, arg := range args {
				if _, err := io.WriteString(writer, arg.Inspect()); err != nil {
					return newError("println: %s", err)
				}
			}

			if _, err := io.WriteString(writer, "\n"); err != nil {
				return newError("println: %s", err)
			}

			return NULL
		},
	}
}
