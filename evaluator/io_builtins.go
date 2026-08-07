package evaluator

import (
	"io"
	"monkey/object"
	"monkey/typing"
)

func init() {
	builtins["puts"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"puts",
				args,
				typing.MinimumArgs(1),
			); err != nil {
				return newError("%s", err.Error())
			}

			stdout, ok := env.Get("STDOUT")
			if !ok {
				return newError("STDOUT is not defined")
			}

			resource, ok := stdout.Value.(*object.Resource)
			if !ok {
				return newError("puts: STDOUT is not a resource")
			}

			writer, ok := resource.Handle.(io.Writer)
			if !ok {
				return newError("puts: STDOUT is not writable")
			}

			for _, arg := range args {
				if _, err := io.WriteString(writer, arg.Inspect()); err != nil {
					return newError("puts: %s", err)
				}
			}

			if _, err := io.WriteString(writer, "\n"); err != nil {
				return newError("puts: %s", err)
			}

			return NULL
		},
	}
}
