package interpreter

import (
	"io"

	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
)

const MaxStackDepth = 64

type State struct {
	Stack        [64]quantity.Q
	StackPointer int

	Units        []quantity.U
	DerivedUnits quantity.UDerivedList

	Output IOutput
	Input  IInput
}

func (s *State) HandleTokensFromInput() (err error) {
	for {
		tok, err := s.Input.ReadToken()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return s.Output.PrintError(err)
		}

		err = s.handleToken(tok)
		if err != nil {
			return s.Output.PrintError(err)
		}
	}
}

func (s *State) handleToken(t syntax.Token) (err error) {
	switch t := t.(type) {
	case *syntax.TokenNumeric:
		return s.LiteralNumber(*t)
	case *syntax.TokenOperator:
		switch t.Literal {
		case "+":
			return s.OperatorPlus()
		case "-":
			return s.OperatorMinus()
		case "*":
			return s.OperatorMultiply()
		case "/":
			return s.OperatorDivide()
		case "p":
			return s.OperatorP()
		case "n":
			return s.OperatorN()
		case "v":
			return s.OperatorV()
		case "c":
			return s.OperatorC()
		case "f":
			return s.OperatorF()
		case "d":
			return s.OperatorD()
		case "r":
			return s.OperatorR()
		default:
			return ErrUnknownOperation{t}
		}
	case *syntax.TokenUnit:
		return s.OperatorUnitConvert(*t)
	default:
		return ErrUnknownOperation{t}
	}
}

func NewDefaultState(input IInput, output IOutput) *State {
	return &State{
		StackPointer: -1,
		Units: []quantity.U{
			quantity.UnitGram,
			quantity.UnitLiter,
			quantity.UnitIU,
			quantity.UnitMeter,
			quantity.UnitMole,
		},
		DerivedUnits: func() (res quantity.UDerivedList) {
			res = append(res, quantity.UnitDerivedAmu)
			res = append(res, quantity.UnitDerivedAmuEng...)
			res = append(res, quantity.UnitDerivedMolar)
			res = append(res, quantity.UnitDerivedMolarEng...)
			res = append(res, quantity.UnitDerivedGramEng...)
			res = append(res, quantity.UnitDerivedLiterEng...)
			res = append(res, quantity.UnitDerivedMeterEng...)
			res = append(res, quantity.UnitDerivedMoleEng...)
			return
		}(),
		Input:  input,
		Output: output,
	}
}
