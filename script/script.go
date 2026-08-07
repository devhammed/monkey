package script

import (
	"fmt"
	"monkey/evaluator"
	"monkey/object"
	"os"
	"path/filepath"
)

// Start runs the script file passed
func Start() {
	file := os.Args[1]
	abs, err := filepath.Abs(file)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Run(string(data), abs, filepath.Dir(abs), true, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}
