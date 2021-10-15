package main

import (
	"bufio"
	"io"
	"os"

	"github.com/eternal-flame-ad/unitdc/interpreter"
	"github.com/eternal-flame-ad/unitdc/repl"
)

var (
	input       = bufio.NewScanner(os.Stdin)
	output      = os.Stdout
	outputError = os.Stderr
)

func main() {
	r := &repl.R{
		Input:     input,
		Output:    output,
		OutputErr: outputError,
	}
	interp := interpreter.NewDefaultState(r, r)
	for {
		if err := r.WritePrompt(); err != nil {
			panic(err)
		}
		if err := r.ParseLineInput(); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if err := interp.HandleTokensFromInput(); err != nil {
			panic(err)
		}
	}
}
