package repl

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"xlang/compiler"
	"xlang/lexer"
	"xlang/parser"
	"xlang/vm"
)

// StartVM starts the REPL with the VM version of xlang
func StartVM(in io.Reader, out io.Writer) {
	fmt.Println(`
	.::        .::.::::::::.::          .::       .::::     .::       .::.:::::::: .::: .::::::    .::::      .::      .::.::            .:       .:::     .::   .::::   
	.::        .::.::      .::       .::   .::  .::    .::  .: .::   .:::.::            .::      .::    .::    .::   .::  .::           .: ::     .: .::   .:: .:    .:: 
	.::   .:   .::.::      .::      .::       .::        .::.:: .:: . .::.::            .::    .::        .::   .:: .::   .::          .:  .::    .:: .::  .::.::        
	.::  .::   .::.::::::  .::      .::       .::        .::.::  .::  .::.::::::        .::    .::        .::     .::     .::         .::   .::   .::  .:: .::.::        
	.:: .: .:: .::.::      .::      .::       .::        .::.::   .:  .::.::            .::    .::        .::   .:: .::   .::        .:::::: .::  .::   .: .::.::   .::::
	.: .:    .::::.::      .::       .::   .::  .::     .:: .::       .::.::            .::      .::     .::   .::   .::  .::       .::       .:: .::    .: :: .::    .: 
	.::        .::.::::::::.::::::::   .::::      .::::     .::       .::.::::::::      .::        .::::      .::      .::.::::::::.::         .::.::      .::  .:::::   

	Made by Gabriel Villalonga in Golang. Followed a book and made research to make an interpreter.
	`)
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			io.WriteString(out, ERROR_MSG+" Error, check them below! "+ERROR_MSG+"\n")
			for n, e := range p.Errors() {
				io.WriteString(out, "\t#"+strconv.Itoa(n)+" "+e+"\n")
			}
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := machine.LastPoppedStackElem()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
