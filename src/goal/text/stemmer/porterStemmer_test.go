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
	stemmer.reset()
	fmt.Println(stemmer)

	//check stem
	stemmer.Stem(dogs)
	fmt.Println(stemmer)

}
