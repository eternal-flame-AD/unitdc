package tokenizer

import (
	"bytes"
	"testing"

	"github.com/eternal-flame-ad/unitdc/syntax"
	. "github.com/smartystreets/goconvey/convey"
)

type tokenizerTestCase struct {
	Source string
	Expect []syntax.Token
}

func TestTokenizer(t *testing.T) {
	Convey("Test Tokenizer", t, func() {
		cases := []tokenizerTestCase{
			{
				Source: "1",
				Expect: []syntax.Token{
					&syntax.TokenNumeric{Literal: "1"},
				},
			},
			{
				Source: " 1\t \n　\t",
				Expect: []syntax.Token{
					&syntax.TokenNumeric{Literal: "1"},
				},
			},
			{
				Source: "1 (ng) \n1 (ul) / p",
				Expect: []syntax.Token{
					&syntax.TokenNumeric{Literal: "1"},
					&syntax.TokenUnit{Literal: "(ng)"},
					&syntax.TokenNumeric{Literal: "1"},
					&syntax.TokenUnit{Literal: "(ul)"},
					&syntax.TokenOperator{Literal: "/"},
					&syntax.TokenOperator{Literal: "p"},
				},
			},
		}
		for _, c := range cases {
			res, err := ParseTokenUntilEOF(bytes.NewBufferString(c.Source))
			So(err, ShouldBeNil)
			So(res, ShouldResemble, c.Expect)
		}
	})
}
