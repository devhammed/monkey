package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/object"
	"os/user"
)

// PROMPT is the repl input prompt
const PROMPT = ">> "

// Start is the repl loop function
func Start(in io.Reader, out io.Writer) {
	user, err := user.Current()
	env := object.NewEnvironment()
	scanner := bufio.NewScanner(in)

	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Hello %s! This is the Monkey programming language!\n",
		user.Username,
	)

	fmt.Printf("Feel free to type in commands.\n")

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		evaluated := evaluator.Run(line, "__REPL__", evaluator.TRUE, env, out)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
