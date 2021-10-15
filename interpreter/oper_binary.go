package interpreter

import "github.com/eternal-flame-ad/unitdc/quantity"

// OperatorPlus pops two quantities from the stack, adds them together
//
// operands must be of equal unit
// result is the same unit as the operand
func (s *State) OperatorPlus() (err error) {
	var operand1, operand2 *quantity.Q
	operand2, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand2)
		}
	}()
	operand1, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand1)
		}
	}()

	if !operand1.UnitExponents.Equal(&operand2.UnitExponents) {
		err = ErrIncompatibleUnit
		return
	}

	res := quantity.Q{
		Number:            operand1.Number + operand2.Number,
		UnitExponents:     operand1.UnitExponents,
		DerivedUnitsToUse: quantity.CombineUDerivedLists(operand1.DerivedUnitsToUse, operand2.DerivedUnitsToUse),
	}
	if operand1.UnitExponents.IsNoUnit() {
		res.UnitExponents = operand2.UnitExponents
	}
	res.UnitExponents.Simplify()
	s.StackPush(res)
	return
}

// OperatorMinus pops two quantities from the stack, subtract the top quantity from the second-to-top quantity
// and push the result onto the stack
//
// operands must be of equal unit
// result is the same unit as the operand
func (s *State) OperatorMinus() (err error) {
	var operand1, operand2 *quantity.Q
	operand2, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand2)
		}
	}()
	operand1, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand1)
		}
	}()

	if !operand1.UnitExponents.Equal(&operand2.UnitExponents) {
		err = ErrIncompatibleUnit
		return
	}

	res := quantity.Q{
		Number:            operand1.Number - operand2.Number,
		UnitExponents:     operand1.UnitExponents,
		DerivedUnitsToUse: quantity.CombineUDerivedLists(operand1.DerivedUnitsToUse, operand2.DerivedUnitsToUse),
	}
	if operand1.UnitExponents.IsNoUnit() {
		res.UnitExponents = operand2.UnitExponents
	}
	res.UnitExponents.Simplify()
	s.StackPush(res)
	return
}

// OperatorMultiply pops two quantities from the stack, multiplies the top
// and push the result onto the stack
//
// unit will be handled accordingly
func (s *State) OperatorMultiply() (err error) {
	var operand1, operand2 *quantity.Q
	operand2, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand2)
		}
	}()
	operand1, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand1)
		}
	}()

	res := quantity.Q{
		Number:            operand1.Number * operand2.Number,
		UnitExponents:     append(operand1.UnitExponents, operand2.UnitExponents...),
		DerivedUnitsToUse: quantity.CombineUDerivedLists(operand1.DerivedUnitsToUse, operand2.DerivedUnitsToUse),
	}
	res.UnitExponents.Simplify()
	s.StackPush(res)
	return
}

// OperatorDivide pops two quantities from the stack, divides them
// and push the result onto the stack
//
// unit will be handled accordingly
func (s *State) OperatorDivide() (err error) {
	var operand1, operand2 *quantity.Q
	operand2, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand2)
		}
	}()
	operand1, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand1)
		}
	}()

	operand2.UnitExponents.Inverse()
	res := quantity.Q{
		Number:            operand1.Number / operand2.Number,
		UnitExponents:     append(operand1.UnitExponents, operand2.UnitExponents...),
		DerivedUnitsToUse: quantity.CombineUDerivedLists(operand1.DerivedUnitsToUse, operand2.DerivedUnitsToUse),
	}
	res.UnitExponents.Simplify()
	s.StackPush(res)
	return
}

// OperatorR reverses the order of the top-most 2 elements on the stack
func (s *State) OperatorR() (err error) {
	var operand1, operand2 *quantity.Q
	operand2, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand2)
		}
	}()
	operand1, err = s.StackPop()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.StackPush(*operand1)
		}
	}()

	s.StackPush(*operand2)
	s.StackPush(*operand1)
	return
}
