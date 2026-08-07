package repl

import (
	"bufio"
	"fmt"
	"monkey/evaluator"
	"monkey/object"
	"os"
	"os/user"
)

// Start is the repl loop function
func Start() {
	currentUser, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Hello %s! This is the Monkey programming language!\n",
		currentUser.Username,
	)

	fmt.Println("Feel free to type in commands.")

	env := object.NewEnvironment()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		evaluated := evaluator.Run(line, "__REPL__", evaluator.TRUE, env)

		if evaluated != nil {
			fmt.Println(evaluated.Inspect())
		}
	}
}
