package tokenize

import (
	"goal/util"
	"strings"
)

type Tokenizer interface {
	//split string into atomic parts
	Tokenize(str string, delimiter string) []string

	//find the boundary for each atomic parts
	TokenizePos(str string, delimiter string) []util.Span
}

func Tokenize(str string, delimiter string) []string {
	return util.SpansToStrings(TokenizePos(str, delimiter), str)
}

func TokenizePos(str string, delimiter string) []*util.Span {
	tokStart := int64(-1)
	var tokens []*util.Span
	inTok := false

	//gather tokens
	for i, v := range str {
		if strings.EqualFold(string(v), delimiter) {
			if inTok {
				tokens = append(tokens, util.NewSpan(tokStart, int64(i)))
				inTok = false
				tokStart = int64(-1)
			}
		} else {
			if !inTok {
				tokStart = int64(i)
				inTok = true
			}
		}
	}

	if inTok {
		tokens = append(tokens, util.NewSpan(tokStart, int64(len(str))))
	}

	return tokens
}
