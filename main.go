package main

import (
	"monkey/repl"
	"monkey/script"
	"os"
)

func main() {
	// Run the script file interpreter if a path is passed else REPL

	if len(os.Args) > 1 {
		script.Start(os.Stdout, os.Args[1:])
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
