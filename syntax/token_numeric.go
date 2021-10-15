package syntax

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TokenNumeric struct {
	Literal string
}

func (n *TokenNumeric) String() string {
	return n.Literal
}

func (n *TokenNumeric) Float() (float64, error) {

	literal := strings.ReplaceAll(n.Literal, "_", "")
	t := strings.SplitN(literal, "e", 3)
	if len(t) == 3 {
		return 0, fmt.Errorf(
			"could not interpret numeric literal %s: more than one exponent notation",
			n.Literal)
	}

	tenExponent := 0
	if len(t) == 2 {
		var err error
		tenExponent, err = strconv.Atoi(t[1])
		if err != nil {
			return 0, fmt.Errorf(
				"could not interpret numeric literal %s: exponent %s is not valid",
				n.Literal,
				t[1],
			)
		}
	}

	numBase, err := strconv.ParseFloat(t[0], 64)
	if err != nil {
		return 0, fmt.Errorf(
			"could not interpret numeric literal %s: base %s is not valid",
			n.Literal,
			t[0],
		)
	}

	return math.Pow10(tenExponent) * numBase, nil
}
