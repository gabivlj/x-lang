package repl

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"xlang/compiler"
	"xlang/lexer"
	"xlang/object"
	"xlang/parser"
	"xlang/vm"
)

const standardLibraryVM = `
let reduce = fn(arr, initial, f) {
	 let iter = fn(arr, result) {
		  if (len(arr) == 0) {
				 return result 
			}
			iter(shift(arr), f(result, first(arr)));
		}
	 iter(arr, initial) 
}

let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), push(accumulated, f(first(arr))));
    }
  };

  iter(arr, []);
};
`

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
	globals := make([]object.Object, vm.GlobalsSize)
	constants := []object.Object{}
	addedStandard := false
	var currentSymbolTable *compiler.SymbolTable
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if !addedStandard {
			line = "\n" + standardLibrary + "\n" + line
			addedStandard = true
		}
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
		var comp *compiler.Compiler
		if currentSymbolTable == nil {
			comp = compiler.New()
		} else {
			comp = compiler.NewWithState(currentSymbolTable, constants)
		}
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}
		bytecode := comp.Bytecode()
		fmt.Println(bytecode.Instructions)
		machine := vm.NewWithGlobalsStore(bytecode, globals)
		constants = bytecode.Constants
		currentSymbolTable = bytecode.Table
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
