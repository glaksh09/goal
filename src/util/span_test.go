package util

import (
	"fmt"
	"testing"
)

func Test_SpanStruct(t *testing.T) {
	//struct
	s1 := Span{1, 4, "hi"}
	fmt.Println("s1:", s1.ToString())

	//constructor
	s2 := NewSpan(2, 6)
	fmt.Println("s2:", s2.ToString())

	//add class
	s2.AddClass("hi")
	fmt.Println("add class to s2:", s2.ToString())

	//add offset
	s2.AddOffset(5)
	fmt.Println("add offset to s2:", s2.ToString())

	//get start/end/class
	fmt.Println("Start:", s2.GetStart())
	fmt.Println("End:", s2.GetEnd())
	fmt.Println("Class:", s2.GetClass())

	//contains Span
	fmt.Println("s2 contains s1:", s2.ContainsSpan(s1))
	s3 := NewSpan(8, 10)
	fmt.Println("s2 contains s3:", s2.ContainsSpan(*s3))

	//contains index
	fmt.Println("s2 contains 5:", s2.ContainsIndex(5))
	fmt.Println("s2 contains 10:", s2.ContainsIndex(10))

	//spans to string
	var spanTok []*Span
	spanTok = append(spanTok, &s1)
	spanTok = append(spanTok, s2)
	spanToStr := SpansToStrings(spanTok, "Hi there, this is a test string")
	fmt.Println("printing strings from spans")
	for _, v := range spanToStr {
		fmt.Println(v)
	}
}
