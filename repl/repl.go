package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
	"github.com/eternal-flame-ad/unitdc/tokenizer"
)

type R struct {
	inputCount  uint64
	outputCount uint64
	Input       *bufio.Scanner
	Output      io.Writer
	OutputErr   io.Writer

	tokenBuf []syntax.Token
}

func (r *R) WritePrompt() error {
	_, err := fmt.Fprintf(r.Output, "In(%d): ", r.inputCount)
	if err != nil {
		return err
	}
	r.inputCount++
	return nil
}

func (r *R) ParseLineInput() error {
	if !r.Input.Scan() {
		return io.EOF
	}
	line := r.Input.Text()
	tok, err := tokenizer.ParseTokenUntilEOF(bytes.NewBufferString(line))
	r.tokenBuf = tok
	if err != nil {
		return r.PrintError(err)
	}
	return nil
}

func (r *R) ReadToken() (tok syntax.Token, err error) {
	if len(r.tokenBuf) == 0 {
		return nil, io.EOF
	}
	tok = r.tokenBuf[0]
	r.tokenBuf = r.tokenBuf[1:]
	return
}

func (r *R) PrintQuantity(values []quantity.Q) (err error) {
	outputPromptHeader := fmt.Sprintf("Out(%d): ", r.outputCount)
	_, err = fmt.Fprint(r.Output, outputPromptHeader, "\n")
	if err != nil {
		return
	}
	r.outputCount++
	for i, value := range values {
		num, unit := value.Format()
		unitStr := ""
		for _, u := range unit {
			unitStr += fmt.Sprintf("(%s)", u.Identifier)
			if u.Exponent != 1 {
				unitStr += strconv.Itoa(u.Exponent) + " "
			}
		}
		_, err = fmt.Fprintf(r.Output, "\t[%d]%f %s\n", i, num, unitStr)
		if err != nil {
			return
		}
	}
	return err
}
func (r *R) PrintError(err error) error {
	output := r.Output
	if r.OutputErr != nil {
		output = r.OutputErr
	}
	_, errOut := fmt.Fprintf(output, "Error: %v\n", err)
	return errOut
}
