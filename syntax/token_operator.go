package syntax

type TokenOperator struct {
	Literal string
}

func (o *TokenOperator) String() string {
	return o.Literal
}
