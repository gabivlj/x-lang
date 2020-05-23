package main

import (
	"fmt"
	"xlang/runtime"
)

func main() {
	// repl.Start(os.Stdin, os.Stdout)
	output, err := runtime.OpenFileAndParse("examples/arrays.xlang")
	if err != nil {
		fmt.Println(err)
	}
	output.Print()
}
