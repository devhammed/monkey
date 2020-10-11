package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"monkey/evaluator"
	"monkey/object"
	"os"
)

// Start runs the script file passed
func Start(out io.Writer, args []string) {
	file := args[0]
	data, err := ioutil.ReadFile(file)

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

	evaluated := evaluator.Run(string(data), file, evaluator.TRUE, env, out)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
		os.Exit(1)
	}
}
