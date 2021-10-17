package interpreter

import "fmt"

type InterpreterError struct {
	UnderlyingError error
}

func (e InterpreterError) Error() string {
	return fmt.Sprintf("invalid operation: %v", e.UnderlyingError)
}

var (
	ErrEmptyStack       = InterpreterError{fmt.Errorf("stack empty")}
	ErrIncompatibleUnit = InterpreterError{fmt.Errorf("incompatible units")}
	ErrUnknownOperation = InterpreterError{fmt.Errorf("unknown operation")}
	ErrUnknownUnit      = InterpreterError{fmt.Errorf("unknown unit identifier")}
)
