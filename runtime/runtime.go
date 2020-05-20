package runtime

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"xlang/eval"
	"xlang/lexer"
	"xlang/object"
	"xlang/parser"
)

const ERROR_MSG = `✖ ✗ ✘ ẋ ☠ ẍ x Ẍ`

// Message is a message of Xlang
type Message struct {
	Line    uint64
	Message []string
}

// Prettify prettifies a message
func (m *Message) Prettify(showLines bool) string {
	s := strings.Builder{}

	for idx, msg := range m.Message {
		if showLines {
			s.WriteString(fmt.Sprintf("\t%d-%d %s\n", m.Line, idx+1, msg))
		} else {
			s.WriteString(fmt.Sprintf("\t%s\n", msg))
		}
	}

	return s.String()
}

// Output output of the program
type Output struct {
	ParseError Message
	Error      Message
	Output     []Message
}

// Print the output of the program
func (o *Output) Print() {
	logMsg := strings.Builder{}
	for _, msg := range o.Output {
		logMsg.WriteString(msg.Prettify(true))
	}
	log.Printf("\nParsing errors: %d\n%sNumber Of Errors: %d\n%sOutput:\n%s", len(o.ParseError.Message), o.ParseError.Prettify(true), len(o.Error.Message), o.Error.Prettify(true), logMsg.String())
}

// OpenFileAndParse parses the program
func OpenFileAndParse(filePath string) (*Output, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	byteCode := make([]byte, stat.Size())
	_, err = file.Read(byteCode)
	if err != nil {
		return nil, err
	}
	code := string(byteCode)
	output := Parse(code)
	return output, nil
}

// Parse .
func Parse(code string) *Output {
	eval := eval.NewEval()
	AddToStandardFunctions(eval)
	output := Output{}

	l := lexer.New(code)
	parser := parser.New(l)
	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		return &Output{ParseError: Message{Line: uint64(program.Line()), Message: parser.Errors()}}
	}
	message := eval.Eval(program)
	if message == nil {
		return &output
	}
	for _, message := range eval.Log {
		if message.Type() == object.LogObject {
			message := message.(*object.Log)
			output.Output = append(output.Output, Message{Line: uint64(message.Line), Message: []string{message.Inspect()}})
		}
	}
	if message.Type() == object.ErrorObject {
		output.Error = Message{Line: uint64(program.Line()), Message: []string{message.Inspect()}}
		return &output
	}

	return &output
}

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

const standardLibrary = `
let reduce = fn(arr, initial, f) {
	 let iter = fn(arr, result) {
		  if (len(arr) == 0) {
				 return result 
			}
			iter(shift(arr), f(result, first(arr)));
		}
	 iter(arr, initial) 
}
`

// AddToStandardFunctions adds util functions into the language
func AddToStandardFunctions(ev *eval.Evaluator) {
	l := lexer.New(standardLibrary)
	p := parser.New(l)
	program := p.ParseProgram()
	ev.Eval(program)
}
