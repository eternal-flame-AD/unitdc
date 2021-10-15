package interpreter

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"

	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
)

type MockedInterpreterInput struct {
	inputTokens []syntax.Token
}

func (i *MockedInterpreterInput) ReadToken() (tok syntax.Token, err error) {
	if len(i.inputTokens) == 0 {
		return nil, io.EOF
	}
	tok = i.inputTokens[0]
	i.inputTokens = i.inputTokens[1:]
	return
}

type MockedInterpreterOutput struct {
	outputErrors     []error
	outputQuantities []quantity.Q
}

func (o *MockedInterpreterOutput) PrintQuantity(values []quantity.Q) (err error) {
	o.outputQuantities = append(o.outputQuantities, values...)
	return nil
}

func (o *MockedInterpreterOutput) PrintError(err error) error {
	o.outputErrors = append(o.outputErrors, err)
	return nil
}

func ShouldExpectOutputQuantities(output interface{}, expected ...interface{}) string {
	out := output.(*MockedInterpreterOutput)
	if len(out.outputQuantities) != len(expected) {
		return fmt.Sprintf("expected %d quantites, but got %d (%#v)", len(expected), len(out.outputQuantities), out.outputQuantities)
	} else {
		for i, e := range expected {
			actual := out.outputQuantities[i]
			expect := e.(quantity.Q)
			if math.Abs(expect.Number-actual.Number) > .001 ||
				!expect.UnitExponents.Equal(&actual.UnitExponents) ||
				!quantity.EqualUDerivedLists(expect.DerivedUnitsToUse, actual.DerivedUnitsToUse) {
				return fmt.Sprintf("[error #%d]: expected quantity %#v, got %#v", i, expect, actual)
			}
		}
	}
	return ""
}

func ShouldExpectOutputErrors(output interface{}, expected ...interface{}) string {
	out := output.(*MockedInterpreterOutput)
	if len(out.outputErrors) != len(expected) {
		return fmt.Sprintf("expected %d errors, but got %d (%#v)", len(expected), len(out.outputErrors), out.outputErrors)
	} else {
		for i, e := range expected {
			switch e := e.(type) {
			case string:
				if !strings.Contains(out.outputErrors[i].Error(), e) {
					return fmt.Sprintf("[error #%d]: expected error string %s, got %#v", i, e, out.outputErrors[i])
				}
			case error:
				if reflect.TypeOf(e) != reflect.TypeOf(out.outputErrors[i]) {
					return fmt.Sprintf("[error #%d]: expected error type %T, got %#v", i, e, out.outputErrors[i])
				}
			}
		}
		out.outputErrors = nil
	}
	return ""
}
