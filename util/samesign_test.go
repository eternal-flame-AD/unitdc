package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSameSign(t *testing.T) {
	Convey("Test SameSign()", t, func() {
		So(SameSign(1, 2), ShouldBeTrue)
		So(SameSign(-1, -2), ShouldBeTrue)
		So(SameSign(0, 1), ShouldBeFalse)
		So(SameSign(-1, 0), ShouldBeFalse)
	})
}
