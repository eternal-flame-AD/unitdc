package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"syscall/js"

	"github.com/eternal-flame-ad/unitdc/interpreter"
	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
	"github.com/eternal-flame-ad/unitdc/tokenizer"
)

type wasmIO struct {
	inputTokens []syntax.Token

	outputFunc js.Value
}

func (w *wasmIO) ReadToken() (tok syntax.Token, err error) {
	if len(w.inputTokens) == 0 {
		return nil, io.EOF
	}
	tok = w.inputTokens[0]
	w.inputTokens = w.inputTokens[1:]
	return
}

func (w *wasmIO) PrintQuantity(values []quantity.Q) (err error) {
	valueStrings := make([]interface{}, len(values))
	for i := range values {
		num, unit := values[i].Format()
		unitStr := ""
		for _, u := range unit {
			unitStr += fmt.Sprintf("(%s)", u.Identifier)
			if u.Exponent != 1 {
				unitStr += strconv.Itoa(u.Exponent) + " "
			}
		}
		valueStrings[i] = fmt.Sprintf("%f %s", num, unitStr)
	}
	w.outputFunc.Invoke(
		"quantity_str",
		valueStrings,
	)
	return nil
}

func (w *wasmIO) PrintError(err error) error {
	w.outputFunc.Invoke(
		"error",
		err.Error(),
	)
	return nil
}

func (w *wasmIO) RequestMoreInput() {
	w.outputFunc.Invoke(
		"ready",
	)
}

func main() {
	wasmio := &wasmIO{}
	interp := interpreter.NewDefaultState(wasmio, wasmio)

	wasmio.outputFunc = js.Global().Get("unitdc_init").Invoke(
		js.FuncOf(func(this js.Value, p []js.Value) interface{} {
			inputType := p[0].String()
			switch inputType {
			case "eval":
				evalDef := p[1]
				code := evalDef.Get("code").String()
				tokens, err := tokenizer.ParseTokenUntilEOF(bytes.NewBufferString(code))
				if err != nil {
					if err := wasmio.PrintError(err); err != nil {
						return js.ValueOf(err.Error())
					}
				}
				wasmio.inputTokens = tokens
				if err := interp.HandleTokensFromInput(); err != nil {
					return js.ValueOf(err.Error())
				}
				wasmio.RequestMoreInput()
			default:
				wasmio.PrintError(fmt.Errorf("unknown WASM ABI input type: %s", inputType))
			}
			return nil
		}),
	)

	wasmio.RequestMoreInput()

	select {}
}
