package main

import (
	"fmt"
	"os"

	"github.com/adrianplavka/fe/repl"
)

func main() {
	fmt.Println("REPL for Fe programming language.\n")
	repl.Start(os.Stdin, os.Stdout)
}
