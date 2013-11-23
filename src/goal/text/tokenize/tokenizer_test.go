package tokenize

import (
	"fmt"
	"testing"
)

func TestTokenizer(t *testing.T) {
	s1 := "hi there, this is a test string"

	//check tokenization
	tokenSpan := TokenizePos(s1, ",")
	fmt.Println("printing spans from string")
	for _, v := range tokenSpan {
		fmt.Println(v.ToString())
	}
	//check tokens
	tokenComma := Tokenize(s1, ",")
	fmt.Println("printing tokens from spans")
	for _, v := range tokenComma {
		fmt.Println(v)
	}

	//check tokens
	tokenSpace := Tokenize(s1, " ")
	fmt.Println("printing tokens from spans")
	for _, v := range tokenSpace {
		fmt.Println(v)
	}

}
