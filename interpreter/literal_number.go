package interpreter

import (
	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
)

// LiteralNumber pushes the number onto the current stack as a unitless quantity.
func (s *State) LiteralNumber(numTok syntax.TokenNumeric) (err error) {
	var num float64
	num, err = numTok.Float()
	if err != nil {
		return
	}
	s.StackPush(quantity.Q{
		Number: num,
	})
	return
}
