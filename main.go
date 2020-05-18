package main

import (
	"os"
	"xlang/repl"
)

func main() {

	repl.Start(os.Stdin, os.Stdout)

}
