package main

import (
	"os"
	"xlang/repl"
)

func main() {
	// http.RunServer()

	// repl.StartVM(os.Stdin, os.Stdout)
	repl.Start(os.Stdin, os.Stdout)
}
