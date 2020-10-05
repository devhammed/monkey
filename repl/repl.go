package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

// PROMPT is the repl input prompt
const PROMPT = ">> "

// Start is the repl loop function
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
