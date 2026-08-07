package script

import (
	"fmt"
	"monkey/evaluator"
	"monkey/object"
	"os"
)

// Start runs the script file passed
func Start() {
	file := os.Args[2]
	data, err := os.ReadFile(file)

	if err != nil {
		fmt.Println("File reading error", err)

		return
	}

	env := object.NewEnvironment()

	evaluated := evaluator.Run(string(data), file, evaluator.TRUE, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}
