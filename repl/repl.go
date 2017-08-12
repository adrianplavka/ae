package repl

import (
	"bufio"
	"fmt"
	"io"

	"../lexer"
	"../token"
)

const prompt = ">> "

// Start the REPL for Fe programming language.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// Loop.
	for {
		// Read.
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Eval.
		line := scanner.Text()
		lex := lexer.New(line)

		// Print.
		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
