package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"monkey/evaluator"
	"monkey/object"
)

// Start runs the script file passed
func Start(out io.Writer, args []string) {
	data, err := ioutil.ReadFile(args[0])

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

	evaluator.Run(string(data), env, out)
}
