package syntax

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type numericTokenTestCase struct {
	Literal string
	Expect  float64
}

func TestNumericToken(t *testing.T) {
	Convey("Numeric Parsing", t, func() {
		Convey("Should be a token", func() {
			tok := &TokenNumeric{}
			var t Token
			t = tok
			_ = t
		})
		Convey("Should parse correctly", func() {
			cases := []numericTokenTestCase{
				{"1", 1},
				{"1.", 1},
				{".1", .1},
				{"0.1", .1},
				{"1_000_000", 1_000_000},
				{"-1", -1},
				{"-.1", -.1},
				{"-.000_1", -.0001},
				{"1e3", 1_000},
				{"1e-3", 0.001},
			}
			for _, e := range cases {
				tok := TokenNumeric{Literal: e.Literal}
				res, err := tok.Float()
				So(err, ShouldBeNil)
				So(res, ShouldAlmostEqual, e.Expect)
			}
		})
	})
}
