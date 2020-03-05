package main

import (
	"fmt"
	"xlang/ast"
	"xlang/testing"
	"xlang/token"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Hello %s! This is XLANG!\n",
	// 	user.Username)
	// fmt.Printf("Being made following a book! Golang rules!\n")
	// repl.Start(os.Stdin, os.Stdout)
	testing.TestParser("let x = 5; return 5;\n")
	p := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	fmt.Println(p.String())
}
