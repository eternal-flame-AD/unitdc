package interpreter

import (
	"math"

	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
)

// OperatorC clear the stack
func (s *State) OperatorC() (err error) {
	s.StackClear()
	return nil
}

// OperatorP print the quantity on top of stack, without altering the stack
func (s *State) OperatorP() (err error) {
	var operand *quantity.Q
	operand, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		s.StackPush(*operand)
	}()

	return s.Output.PrintQuantity([]quantity.Q{*operand})
}

// OperatorN pop and print the quantity of top of stack
func (s *State) OperatorN() (err error) {
	var operand *quantity.Q
	operand, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand)
		}
	}()

	return s.Output.PrintQuantity([]quantity.Q{*operand})
}

// OperatorV pops the quantity of top of stack, and pushes its square root onto the stack
//
// All involved base units must be of even exponents.
func (s *State) OperatorV() (err error) {
	var operand *quantity.Q
	operand, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand)
		}
	}()

	res := quantity.Q{
		Number:            math.Sqrt(operand.Number),
		DerivedUnitsToUse: operand.DerivedUnitsToUse,
		UnitExponents:     operand.UnitExponents,
	}
	res.UnitExponents.Simplify()
	for i := range res.UnitExponents {
		if res.UnitExponents[i].Exponent%2 != 0 {
			err = ErrIncompatibleUnit
			return
		}
		res.UnitExponents[i].Exponent /= 2
	}
	s.StackPush(res)
	return
}

// OperatorF prints the content of the whole stack
func (s *State) OperatorF() (err error) {
	return s.Output.PrintQuantity(s.Stack[:s.StackPointer+1])
}

// OperatorD pushes a copy of the quantity on top of stack onto the stack
func (s *State) OperatorD() (err error) {
	var operand *quantity.Q
	operand, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand)
		}
	}()

	s.StackPush(*operand)

	operandCopy := *operand
	operandCopy.UnitExponents = operandCopy.UnitExponents.Clone()
	operandCopy.DerivedUnitsToUse = operandCopy.DerivedUnitsToUse.Clone()
	s.StackPush(operandCopy)
	return
}

// OperatorUnit modifies the unit on the top of the stack
//
//
// if the provided unit is no unit (1), the quantity on the top of the stack will be changed to no unit
//
// if the quantity on top of the stack is of no unit, it will be assigned the provided unit
//
//
// **Additionally**, this operator will change the preferred display method of the quantity of the top of the stack:
//
//
// If the provided unit is a base unit:
// All preferred derived units that contain this base will be removed from the preferred list.
//
// Example: 1 (M) (l) will result in a quantity of 1 (M), but displayed as 1 (mol)(l)-1
//
// Example: 1 (ml) (l) will result in a quantity of .001 (l).
//
//
// If the provided unit is a derived, simple unit (the same as a base unit, just a multiplier and/or a offset):
// The provided unit will be added to the preferred list, and
// all existing preferred derived units that share the base unit as this new preferred unit will be removed from the preferred list.
//
// Example: 1 (l) (ml) will result in a quantity of 1 (l). Displayed as 1000 (ml). This expression is equivalent to 1000 (ml).
//
// Example: 10 (ml) (ul) will result in a quantity of .01 (l). Displayed as 10000 (ul).
//
// If the provided unit is a derived, complex unit (consist of more than one base unit, or one base unit of exponent other than 1):
// The provided unit will be added to the preferred list, and
// all existing preferred derived units that share any base unit used in the new preferred unit will be removed from the preferred list.
//
// Example: 1 (l) (M) will result in a quantity of 1 (l). Displayed as 1 (l), because use of molar in this context does not simplify the unit.
//
// Example: 1 (mol) (M) 1 (l) / will result in a quantity of 1 (mol)(l)-1. Displayed as 1 (M).
//
// Example: 1 (mol) 1 (l) / (M) will result in a quantity of 1 (mol)(l)-1. Displayed as 1 (M). This is equivalent to the previous example but more readable.
//
// Example: 1 (mol) 1 (ml) / (uM) will result in a quantity of 1 (mol)(l)-1. Displayed as 1000000 (uM).
//
// Real-life examples on the use of unit operators:
//
//
// Question 1:
//   What is the resulting concentration of 10 grams of NaCl in 10 liters of solution, expressed in uM?
//   10 (g)              # mass of solid used
//   58.44 (g) 1 (mol) / # MW of NaCl 58.44 g/mol
//   /                   # divide the two quantities, result is (mol) of NaCl in solution
//   10 (l) /            # compute the molarity, in (mol)(l)-1
//   (uM) p              # display the result in (uM) unit.
//
// Question 2:
//   What is the mass of a 100uM 10-nt ssDNA in a 10 uL solution, expressed in nanograms?
//
//   303.7 10 * 79 +                    # compute the MW of the ssDNA, based on estimation equation
//                                      # https://www.thermofisher.com/us/en/home/references/ambion-tech-support/rna-tools-and-calculators/dna-and-rna-molecular-weights-and-conversions.html
//   (g) 1 (mol) /                      # set the unit to g/mol
//
//   10 (ul) 100 (uM) *                 # compute the moles of DNA in the solution
//   *                                  # multiple by MW to get the mass of DNA
//   (ng) p                             # display the result in (ng) unit
func (s *State) OperatorUnitConvert(unitTok syntax.TokenUnit) (err error) {

	var unit string
	unit, err = unitTok.UnitIdentifier()
	if err != nil {
		return
	}

	var operand *quantity.Q
	operand, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		s.StackPush(*operand)
	}()

	if unit == "1" {
		operand.UnitExponents = quantity.UCombination{}
		return
	}

	if operand.UnitExponents.IsNoUnit() {
		for _, u := range s.DerivedUnits {
			if u.Identifier == unit {
				operand.UnitExponents = make(quantity.UCombination, len(u.UnitExponents))
				copy(operand.UnitExponents, u.UnitExponents)
				operand.Number = operand.Number*u.Multiplier + u.Offset
			}
		}
		for _, u := range s.Units {
			if u.Identifier == unit {
				operand.UnitExponents = quantity.UCombination{quantity.UExp{Unit: u, Exponent: 1}}
			}
		}
	}
	for _, u := range s.DerivedUnits {
		if u.Identifier == unit {
			var newDerivedUnitsToUse quantity.UDerivedList
			for _, ud := range operand.DerivedUnitsToUse {
				if !ud.UnitExponents.HasOverlap(u.UnitExponents) {
					newDerivedUnitsToUse = append(newDerivedUnitsToUse, ud)
				}
			}
			newDerivedUnitsToUse = append(newDerivedUnitsToUse, u)
			operand.DerivedUnitsToUse = newDerivedUnitsToUse
		}
	}
	for _, u := range s.Units {
		if u.Identifier == unit {
			// clear all preferred derived units the contain this unit
			var newDerivedUnitsToUse quantity.UDerivedList
			for _, ud := range operand.DerivedUnitsToUse {
				conflict := false
				for _, udu := range ud.UnitExponents {
					if udu.Unit.ID == u.ID {
						conflict = true
						break
					}
				}
				if !conflict {
					newDerivedUnitsToUse = append(newDerivedUnitsToUse, ud)
				}
			}
			operand.DerivedUnitsToUse = newDerivedUnitsToUse
		}
	}

	return
}
