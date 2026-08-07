package script

import (
	"fmt"
	"monkey/evaluator"
	"monkey/object"
	"os"
)

// Start runs the script file passed
func Start() {
	args := os.Args[1:]
	file := args[0]
	data, err := os.ReadFile(file)

	if err != nil {
		fmt.Println("File reading error", err)

		return
	}

	env := object.NewEnvironment()

	var scriptArgs []object.Object

	for _, scriptArg := range args {
		scriptArgs = append(scriptArgs, &object.String{Value: scriptArg})
	}

	env.Set("ARGV", &object.Array{Elements: scriptArgs})

	evaluated := evaluator.Run(string(data), file, evaluator.TRUE, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}
