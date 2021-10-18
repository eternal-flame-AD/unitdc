package syntax

import (
	"fmt"
	"math/big"
	"strings"
)

type TokenNumeric struct {
	Literal string
}

func (n *TokenNumeric) String() string {
	return n.Literal
}

func (n *TokenNumeric) BigFloat(f *big.Float) error {
	_, _, err := f.Parse(
		strings.ReplaceAll(n.Literal, "_", ""), 10)
	return err
}

func (n *TokenNumeric) Float() (float64, error) {
	f := big.NewFloat(0)
	if err := n.BigFloat(f); err != nil {
		return 0, fmt.Errorf(
			"could not interpret numeric literal %s: %w",
			n.Literal,
			err)
	}
	ret, _ := f.Float64()
	return ret, nil
}
