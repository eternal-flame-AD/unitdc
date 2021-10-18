package interpreter

import (
	"fmt"

	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
)

type ErrEmptyStack struct {
}

func (e ErrEmptyStack) Error() string {
	return fmt.Sprintf("stack empty")
}

type ErrIncompatibleUnit struct {
	TargetUnit    quantity.UCombination
	OffendingUnit quantity.UCombination
}

func (e ErrIncompatibleUnit) Error() string {
	if e.TargetUnit != nil {
		return fmt.Sprintf("incompatible units: could not coerce %v into %v",
			e.OffendingUnit, e.TargetUnit,
		)
	}
	return fmt.Sprintf("incompatible units: %v is unacceptable for this operation",
		e.OffendingUnit,
	)
}

type ErrUnknownOperation struct {
	Token syntax.Token
}

func (e ErrUnknownOperation) Error() string {
	return fmt.Sprintf("undefined operation: %v", e.Token)
}

type ErrUnknownUnit struct {
	UnitIdentifier string
}

func (e ErrUnknownUnit) Error() string {
	return fmt.Sprintf("undefined unit: (%s)", e.UnitIdentifier)
}
