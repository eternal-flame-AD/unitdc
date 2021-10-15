package tokenizer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode"

	"github.com/eternal-flame-ad/unitdc/syntax"
)

var (
	operatorTokenRegexp = regexp.MustCompile("^[cdrbpnPf+\\-*/]$")
	numericTokenRegexp  = regexp.MustCompile("^(\\+|-)?[0-9._]+(e(\\+|-)?[0-9_]+)?$")
	unitTokenRegexp     = regexp.MustCompile("^\\(1|[a-zA-Z]\\w*\\)$")
)

func isWhiteSpace(c rune) bool {
	return unicode.IsSpace(c)
}

func isNewLine(c rune) bool {
	return c == '\r' || c == '\n'
}

func ParseTokenUntilEOF(r io.RuneReader) (res []syntax.Token, err error) {
	for {
		var tok syntax.Token
		tok, err = ParseToken(r)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		res = append(res, tok)
	}
}

func ParseToken(r io.RuneReader) (syntax.Token, error) {
	// discard white space
	var nextRune rune

	inComment := false
	var err error
	for {
		nextRune, _, err = r.ReadRune()
		if err != nil {
			return nil, err
		}
		if isNewLine(nextRune) {
			inComment = false
		}
		if !isWhiteSpace(nextRune) {
			if nextRune == '#' {
				inComment = true
			}
			if !inComment {
				break
			}
		}
	}

	var tokenBuf bytes.Buffer
	tokenBuf.WriteRune(nextRune)
	for {
		nextRune, _, err = r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if isWhiteSpace(nextRune) {
			break
		}
		tokenBuf.WriteRune(nextRune)
	}

	tokenLiteral := tokenBuf.String()
	if unitTokenRegexp.MatchString(tokenLiteral) {
		return &syntax.TokenUnit{Literal: tokenLiteral}, nil
	} else if numericTokenRegexp.MatchString(tokenLiteral) {
		return &syntax.TokenNumeric{Literal: tokenLiteral}, nil
	} else if operatorTokenRegexp.MatchString(tokenLiteral) {
		return &syntax.TokenOperator{Literal: tokenLiteral}, nil
	}

	return nil, fmt.Errorf("cannot parse token %s", tokenLiteral)
}
