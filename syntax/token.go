package syntax

import (
	"fmt"
)

type Token interface {
	fmt.Stringer
}
