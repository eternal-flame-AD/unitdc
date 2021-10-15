package interpreter

import (
	"testing"

	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInterpreter(t *testing.T) {

	Convey("Test Interpreter", t, func() {
		mockInput := new(MockedInterpreterInput)
		mockOutput := new(MockedInterpreterOutput)
		mockInterpreter := NewDefaultState(mockInput, mockOutput)

		testErrorOnOneOperand := func(operator string) {
			Convey("should error when stack only has one operand", func() {
				for mockInterpreter.StackDepth() > 1 {
					mockInterpreter.StackPop()
				}
				So(mockInterpreter.StackDepth(), ShouldEqual, 1)
				mockInput.inputTokens = []syntax.Token{
					&syntax.TokenOperator{Literal: operator},
				}
				So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
				So(mockOutput, ShouldExpectOutputErrors, ErrEmptyStack)
				Convey("error should not alter stack", func() {
					So(mockInterpreter.StackDepth(), ShouldEqual, 1)
					So(mockInterpreter.Stack[mockInterpreter.StackPointer].Number,
						ShouldAlmostEqual, 1)
				})
			})
		}
		testErrorOnNoOperands := func(operator string) {
			Convey("should error when stack only has no operands", func() {
				mockInterpreter.StackClear()
				mockInput.inputTokens = []syntax.Token{
					&syntax.TokenOperator{Literal: "+"},
				}
				So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
				So(mockOutput, ShouldExpectOutputErrors, ErrEmptyStack)
				Convey("error should not alter stack", func() {
					So(mockInterpreter.StackDepth(), ShouldEqual, 0)
				})
			})
		}

		Convey("Stack Operations", func() {
			Convey("StackClear() should empty stack", func() {
				mockInterpreter.StackPush(quantity.Q{Number: 1})
				mockInterpreter.StackPush(quantity.Q{Number: 2})
				So(mockInterpreter.StackDepth(), ShouldEqual, 2)
				mockInterpreter.StackClear()
				So(mockInterpreter.StackDepth(), ShouldEqual, 0)
			})
			Convey("Empty Stack should return error", func() {
				mockInterpreter.StackPush(quantity.Q{Number: 1})
				mockInterpreter.StackPush(quantity.Q{Number: 2})
				v1, err1 := mockInterpreter.StackPop()
				v2, err2 := mockInterpreter.StackPop()
				v3, err3 := mockInterpreter.StackPop()
				So(v1.Number, ShouldAlmostEqual, 2)
				So(v2.Number, ShouldAlmostEqual, 1)
				So(v3, ShouldBeNil)
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				So(err3, ShouldBeError, ErrEmptyStack)
			})
			Convey("Full stack should leak", func() {
				for i := 0; i <= MaxStackDepth; i++ {
					mockInterpreter.StackPush(quantity.Q{Number: float64(i)})
				}
				So(mockInterpreter.StackDepth(), ShouldEqual, MaxStackDepth)
				for i := MaxStackDepth; i > 0; i-- {
					v, err := mockInterpreter.StackPop()
					So(err, ShouldBeNil)
					So(v.Number, ShouldAlmostEqual, float64(i))
				}
				v, err := mockInterpreter.StackPop()
				So(err, ShouldBeError, ErrEmptyStack)
				So(v, ShouldBeNil)
			})
		})

		Convey("RPN Token Handling", func() {
			Convey("Number literal should push onto stack", func() {
				mockInput.inputTokens = []syntax.Token{
					&syntax.TokenNumeric{Literal: "1"},
					&syntax.TokenNumeric{Literal: "1.5"},
					&syntax.TokenNumeric{Literal: ".25"},
				}
				So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
				So(mockInput.inputTokens, ShouldBeEmpty)
				So(mockInterpreter.StackDepth(), ShouldEqual, 3)
				So(mockOutput, ShouldExpectOutputErrors)
				Convey("I/O Operators", func() {
					Convey("Operator p", func() {
						mockInput.inputTokens = []syntax.Token{
							&syntax.TokenOperator{Literal: "p"},
							&syntax.TokenOperator{Literal: "p"},
						}
						So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
						So(mockInput.inputTokens, ShouldBeEmpty)
						So(mockOutput, ShouldExpectOutputErrors)

						Convey("Should not alter stack", func() {
							So(mockInterpreter.StackDepth(), ShouldEqual, 3)
						})
						Convey("Should output stack top", func() {
							So(mockOutput, ShouldExpectOutputQuantities, quantity.Q{Number: .25}, quantity.Q{Number: .25})
						})
						testErrorOnNoOperands("p")
					})
					Convey("Operator n", func() {
						mockInput.inputTokens = []syntax.Token{
							&syntax.TokenOperator{Literal: "n"},
							&syntax.TokenOperator{Literal: "n"},
						}
						So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
						So(mockInput.inputTokens, ShouldBeEmpty)
						So(mockOutput, ShouldExpectOutputErrors)

						Convey("Should pop stack", func() {
							So(mockInterpreter.StackDepth(), ShouldEqual, 1)
						})
						Convey("Should output stack top", func() {
							So(mockOutput, ShouldExpectOutputQuantities,
								quantity.Q{Number: .25},
								quantity.Q{Number: 1.5})
						})
						testErrorOnNoOperands("n")
					})
					Convey("Operator f", func() {
						mockInput.inputTokens = []syntax.Token{
							&syntax.TokenOperator{Literal: "f"},
							&syntax.TokenOperator{Literal: "f"},
						}
						So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
						So(mockInput.inputTokens, ShouldBeEmpty)
						So(mockOutput, ShouldExpectOutputErrors)

						Convey("Should not alter stack", func() {
							So(mockInterpreter.StackDepth(), ShouldEqual, 3)
						})
						Convey("Should output full stack", func() {
							So(mockOutput, ShouldExpectOutputQuantities,
								quantity.Q{Number: 1},
								quantity.Q{Number: 1.5},
								quantity.Q{Number: .25},
								quantity.Q{Number: 1},
								quantity.Q{Number: 1.5},
								quantity.Q{Number: .25},
							)
						})
						Convey("Should be no-op on empty stack", func() {
							mockInterpreter.StackClear()
							mockInput.inputTokens = []syntax.Token{
								&syntax.TokenOperator{Literal: "f"},
							}
							So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
							So(mockInput.inputTokens, ShouldBeEmpty)
							So(mockOutput, ShouldExpectOutputErrors)
						})
					})
				})
				Convey("Operators on no unit operands", func() {

					Convey("Unit Operator", func() {
						Convey("Unit operator should set unit", func() {
							mockInput.inputTokens = []syntax.Token{
								&syntax.TokenUnit{Literal: "(l)"},
							}
							So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
							So(mockInterpreter.StackDepth(), ShouldEqual, 3)
							So(mockOutput, ShouldExpectOutputErrors)
							So(mockInterpreter.Stack[mockInterpreter.StackPointer].Number,
								ShouldAlmostEqual,
								.25,
							)
							So(mockInterpreter.Stack[mockInterpreter.StackPointer].UnitExponents, ShouldResemble,
								quantity.UCombination{{Unit: quantity.UnitLiter, Exponent: 1}})
						})
					})
					Convey("Operator +", func() {
						Convey("should add to no unit", func() {
							mockInput.inputTokens = []syntax.Token{
								&syntax.TokenOperator{Literal: "+"},
							}
							So(mockInterpreter.HandleTokensFromInput(), ShouldBeNil)
							So(mockInterpreter.StackDepth(), ShouldEqual, 2)
							So(mockOutput, ShouldExpectOutputErrors)
							So(mockInterpreter.Stack[mockInterpreter.StackPointer].Number,
								ShouldAlmostEqual,
								1.5+.25,
							)
							So(mockInterpreter.Stack[mockInterpreter.StackPointer].UnitExponents, ShouldBeEmpty)
						})
						testErrorOnNoOperands("+")
						testErrorOnOneOperand("+")
					})
				})
			})
		})
	})
}
