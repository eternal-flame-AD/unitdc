package syntax

import "strings"

type TokenUnit struct {
	Literal string
}

func (u *TokenUnit) String() string {
	return u.Literal
}

func (u *TokenUnit) UnitIdentifier() (string, error) {
	return strings.TrimSpace(
		strings.TrimPrefix(strings.TrimSuffix(
			u.Literal,
			")",
		), "(",
		),
	), nil
}
