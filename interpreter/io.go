package interpreter

import (
	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
)

type IInput interface {
	ReadToken() (tok syntax.Token, err error)
}

type IOutput interface {
	PrintQuantity(values []quantity.Q) (err error)
	PrintError(err error) error
}
