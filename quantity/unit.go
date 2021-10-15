package quantity

type UExp struct {
	Unit     U
	Exponent int
}

type U struct {
	Identifier string
	ID         int
}
