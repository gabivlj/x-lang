package testing

import (
	"fmt"
	"xlang/lexer"
	"xlang/parser"
)

// TestParser .
func TestParser(s string) {
	p := parser.New(lexer.New(s))
	program := p.ParseProgram()
	fmt.Println(program.String())
}
