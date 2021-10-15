package syntax

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type unitTokenTestCase struct {
	Literal string
	Expect  string
}

func TestUnitToken(t *testing.T) {
	Convey("Unit Parsing", t, func() {
		Convey("Should be a token", func() {
			tok := &TokenUnit{}
			var t Token
			t = tok
			_ = t
		})
		Convey("Should parse correctly", func() {
			cases := []unitTokenTestCase{
				{"(mg)", "mg"},
				{"(\tmcg )", "mcg"},
			}
			for _, e := range cases {
				tok := TokenUnit{e.Literal}
				res, err := tok.UnitIdentifier()
				So(err, ShouldBeNil)
				So(res, ShouldEqual, e.Expect)
			}
		})
	})
}
