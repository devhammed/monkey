package evaluator

import (
	"fmt"
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
				return newError(err.Error())
			}

			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}

			fmt.Println("")

			return NULL
		},
	}
}
