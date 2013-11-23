package stemmer

import (
	"fmt"
	"testing"
)

func TestStemmer(t *testing.T) {
	//s1 := "hi there, this is a test string"

	//check initialization
	stemmer := new(PorterStemmer)
	fmt.Println(stemmer)

	//check reset()
	stemmer.i = 50
	stemmer.dirty = true
	fmt.Println(stemmer)
	stemmer.Reset()
	fmt.Println(stemmer)

	//check stem
	str1 := "dogs"
	stemmer.Stem(str1)
	fmt.Println(str1, " -> ", string(stemmer.b))

}
