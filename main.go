package main

import (
	"monkey/repl"
	"monkey/script"
	"os"
)

// Run the script file interpreter if a path is passed else REPL
func main() {
	if len(os.Args) > 1 {
		script.Start()
	} else {
		repl.Start()
	}
}
