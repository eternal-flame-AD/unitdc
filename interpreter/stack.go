package interpreter

import "github.com/eternal-flame-ad/unitdc/quantity"

func (s *State) StackCopy() (q []quantity.Q) {
	q = make([]quantity.Q, s.StackDepth())
	for i := range q {
		q[i] = quantity.Q{
			Number:            s.Stack[i].Number,
			UnitExponents:     s.Stack[i].UnitExponents.Clone(),
			DerivedUnitsToUse: s.Stack[i].DerivedUnitsToUse.Clone(),
		}
	}
	return
}

func (s *State) StackPush(q quantity.Q) {
	q.UnitExponents = q.UnitExponents.Clone()
	q.DerivedUnitsToUse = q.DerivedUnitsToUse.Clone()
	if s.StackPointer == MaxStackDepth-1 {
		// stack is full
		copy(s.Stack[:], s.Stack[1:])
		s.Stack[s.StackPointer] = q
	} else {
		s.StackPointer++
		s.Stack[s.StackPointer] = q
	}
}

func (s *State) StackPop() (q *quantity.Q, err error) {
	if s.StackPointer == -1 {
		return nil, ErrEmptyStack{}
	}
	q = &quantity.Q{
		Number:            s.Stack[s.StackPointer].Number,
		UnitExponents:     s.Stack[s.StackPointer].UnitExponents.Clone(),
		DerivedUnitsToUse: s.Stack[s.StackPointer].DerivedUnitsToUse.Clone(),
	}
	s.StackPointer--
	return
}

func (s *State) StackClear() {
	s.StackPointer = -1
}

func (s *State) StackDepth() int {
	return s.StackPointer + 1
}
