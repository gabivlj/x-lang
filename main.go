package main

import (
	"xlang/testing"
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
	testing.TestParser("foobar; \n 5; -6;")
	// p := &ast.Program{
	// 	Statements: []ast.Statement{
	// 		&ast.LetStatement{
	// 			Token: token.Token{Type: token.LET, Literal: "let"},
	// 			Name: &ast.Identifier{
	// 				Token: token.Token{Type: token.IDENT, Literal: "myVar"},
	// 				Value: "myVar",
	// 			},
	// 			Value: &ast.Identifier{
	// 				Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
	// 				Value: "anotherVar",
	// 			},
	// 		},
	// 	},
	// }
	// fmt.Println(p.String())
}
