package repl

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"xlang/eval"
	"xlang/lexer"
	"xlang/parser"
)

const PROMPT = ">> "
const ERROR_MSG = `✖ ✗ ✘ ẋ ☠ ẍ x Ẍ`

// Start stars the REPL of Xlang.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
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
	evaluator := eval.NewEval()
	AddToStandardFunctions(evaluator)
	for {
		fmt.Printf(PROMPT)
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
		io.WriteString(out, "Built AST succesfully, running...\n")
		evaluatedProgram := evaluator.Eval(program)
		if evaluatedProgram == nil {
			continue
		}
		io.WriteString(out, evaluatedProgram.Inspect())
		io.WriteString(out, "\n")
	}
}

// let reduce = fn(arr, initial, f) {
// 	let iter = fn(arr, result) {
// 		 if (len(arr) == 0) {
// 				return result
// 		 }
// 		 iter(shift(arr), f(result, first(arr)));
// 	 }
// 	iter(arr, initial)
// }
const standardLibrary = `


let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(shift(arr), push(accumulated, f(first(arr))));
    }
  };

  iter(arr, []);
};

let wrapper = fn() {
	let countDown = fn(x) {
			if (x == 0) {
					return 1;
			} else {
					countDown(x - 1);
			}
	};
	countDown(1);
};
wrapper();
`

// AddToStandardFunctions adds util functions into the language
func AddToStandardFunctions(ev *eval.Evaluator) {
	l := lexer.New(standardLibrary)
	p := parser.New(l)
	program := p.ParseProgram()
	ev.Eval(program)
}
