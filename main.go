package main

import (
	"monkey/repl"
	"monkey/script"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		script.Start(os.Stdout, os.Args[1:])
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
