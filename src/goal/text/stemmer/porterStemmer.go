package stemmer

const INC = 50

type PorterStemmer struct {
	b           []byte
	i, j, k, k0 int
	dirty       bool
}

//reset the stemmer to stem another word
//call reset() if calling add(byte) and then stem()
func (s *PorterStemmer) Reset() error {
	s.b = new([]byte)
	s.i = 0
	s.dirty = false
	return nil
}

//Add a character to the word being stemmed.
//After finishing adding characters, call stem(void)
//to process the word
func (s *PorterStemmer) Add(ch byte) error {

	append(s.b, ch)
	/*
		if len(s.b) == s.i {
			new_b:= new([]byte,s.i+INC)
			for c:=0;c<s.i;c++{
				new_b[c] = s.b[c]
			}
			s.b = new_b
		}
		s.b[s.i++] = ch
	*/
	return nil
}

//convert the stemmed slice of word to string
func (s *PorterStemmer) ToString() (string, error) {
	return string(s.b), nil
}

//get length of the stemmed word
func (s *PorterStemmer) GetResultLength() (int, error) {
	return s.i, nil
}

//get result buffer slice
func (s *PorterStemmer) GetResultBuffer() ([]byte, error) {
	return s.b, nil
}

//check if character at index i is a consonant
func (s *PorterStemmer) IsCons(i int) bool {
	switch s.b[i] {
	case 'a', 'e', 'i', 'o', 'u':
		return false
	case 'y':
		if i == s.k0 {
			return true
		} else {
			return !cons(i - 1)
		}
	default:
		return true
	}
}

//measures the number of consonant sequences
func (s *PorterStemmer) CountConsSeq() (int, error) {
	n := 0
	i := s.k0
	for {
		if i > s.j {
			return n
		}
		if !s.IsCons(i) {
			break
		}
		i++
	}
	i++
	for {
		for {
			if i > s.j {
				return n
			}
			if s.IsCons(i) {
				break
			}
			i++
		}
		i++
		n++
		for {
			if i > s.j {
				return n
			}
			if !s.IsCons(i) {
				break
			}
			i++
		}
		i++
	}
}

//check if there is a vowel in stem
func (s *PorterStemmer) IsVowelInStem() bool {
	for i := s.k0; i <= s.j; i++ {
		if !s.IsCons(i) {
			return true
		}
	}
	return false
}

//check if two consecutive consonants are present
func (s *PorterStemmer) IsConsecutiveCons(j int) bool {
	if j < s.k0 {
		return false
	}
	if s.b[j] != s.b[j-1] {
		return false
	}
	return s.IsCons(j)
}

//count consonant-vowel-consonant pattern
func (s *PorterStemmer) IsCVCPattern(i int) bool {
	if i < s.k0+2 || !s.IsCons(i) || s.IsCons(i-1) || !s.IsCons(i-2) {
		return false
	} else {
		if ch := s.b[i]; ch == byte('w') || ch == byte('x') || ch == byte('y') {
			return false
		}
	}
	return true
}

func (s *PorterStemmer) Ends(str string) bool {
	l := len(str)
	o := s.k - l + 1
	if o < s.k0 {
		return false
	}
	for i := 0; i < l; i++ {
		if s.b[o+i] != str[i] {
			return false
		}
	}
	s.j = s.k - l
	return true
}

//set the buffer to specified string
func (s *PorterStemmer) SetBuffer(str string) error {
	l := len(str)
	o = s.j + 1
	for i := 0; i < l; i++ {
		s.b[o+i] = str[i]
	}
	s.k = s.j + l
	s.dirty = true
	return nil
}

//set buffer based on consonant sequences
func (s *PorterStemmer) SetBufferOnConsSeq(str string) {
	if s.CountConsSeq() > 0 {
		s.SetBuffer(str)
	}
	return nil
}

/*
step 1: remove plurals and -ed or -ing
dogs -> dogs
*/
func (s *PorterStemmer) Step1() error {
	if s.b[s.k] == 's' {
		if s.Ends("sses") {
			s.k = s.k - 2
		} else if s.Ends("ies") {
			s.SetBuffer("i")
		} else if s.b[s.k-1] != 's' {
			s.k = s.k - 1
		}
	}

	if s.Ends("eed") {
		if s.CountConsSeq > 0 {
			s.k = s.k - 1
		}
	} else if (s.Ends("ed") || s.Ends("ing")) && s.IsVowelInStem() {
		s.k = s.j
		if s.Ends("at") {
			s.SetBuffer("ate")
		} else if s.Ends("bl") {
			s.SetBuffer("ble")
		} else if s.Ends("iz") {
			s.SetBuffer("ize")
		} else if s.IsConsecutiveCons(s.k) {
			ch := s.b[s.k]
			s.k--
			if ch == 'l' || ch == 's' || ch == 'z' {
				s.k++
			}
		} else if s.CountConsSeq() == 1 && s.IsCVCPattern(s.k) {
			s.SetBuffer("e")
		}
	}

	return nil
}

/*
Step 2:
change 'y' to 'i' when another vowel is present in the stem
*/
func (s *PorterStemmer) Step2() error {
	if s.Ends("y") && s.IsVowelInStem() {
		s.b[s.k] = 'i'
		s.dirty = true
	}
	return nil
}

/*
Step 3: change double suffices to sigle ones
-ization (-ize and -ation ) -> -ize
Precondition: CountConsSeq(string before the suffice)>0
*/
func (s *PorterStemmer) Step3() error {
	if s.k == s.k0 {
		return nil
	}

	switch s.b[k-1] {
	case 'a':
		if s.Ends("ational") {
			s.SetBufferOnConsSeq("ate")
			break
		}
		if s.Ends("tional") {
			s.SetBufferOnConsSeq("tion")
			break
		}
	case 'c':
		if s.Ends("enci") {
			s.SetBufferOnConsSeq("ence")
			break
		}
		if s.Ends("anci") {
			s.SetBufferOnConsSeq("ance")
			break
		}
	case 'e':
		if s.Ends("izer") {
			s.SetBufferOnConsSeq("ize")
			break
		}
	case 'l':
		if s.Ends("bli") {
			s.SetBufferOnConsSeq("ble")
			break
		}
		if s.Ends("alli") {
			s.SetBufferOnConsSeq("al")
			break
		}
		if s.Ends("entli") {
			s.SetBufferOnConsSeq("ent")
			break
		}
		if s.Ends("eli") {
			s.SetBufferOnConsSeq("e")
			break
		}
		if s.Ends("ousli") {
			s.SetBufferOnConsSeq("ous")
			break
		}
	case 'o':
		if s.Ends("ization") {
			s.SetBufferOnConsSeq("ize")
			break
		}
		if s.Ends("ation") {
			s.SetBufferOnConsSeq("ate")
			break
		}
		if s.Ends("ator") {
			s.SetBufferOnConsSeq("ate")
			break
		}
	case 's':
		if s.Ends("alism") {
			s.SetBufferOnConsSeq("al")
			break
		}
		if s.Ends("iveness") {
			s.SetBufferOnConsSeq("ive")
			break
		}
		if s.Ends("fulness") {
			s.SetBufferOnConsSeq("ful")
			break
		}
		if s.Ends("ousness") {
			s.SetBufferOnConsSeq("ous")
			break
		}
	case 't':
		if s.Ends("aliti") {
			s.SetBufferOnConsSeq("al")
			break
		}
		if s.Ends("iviti") {
			s.SetBufferOnConsSeq("ive")
			break
		}
		if s.Ends("biliti") {
			s.SetBufferOnConsSeq("ble")
			break
		}
	case 'g':
		if s.Ends("logi") {
			s.SetBufferOnConsSeq("log")
			break
		}

	}
	return nil
}

/*
Step 4: handle -ic-, -full, -ness
*/
func (s *PorterStemmer) Step4() error {

	switch s.b[k-1] {
	case 'e':
		if s.Ends("icate") {
			s.SetBufferOnConsSeq("ic")
			break
		}
		if s.Ends("ative") {
			s.SetBufferOnConsSeq("")
			break
		}
		if s.Ends("alize") {
			s.SetBufferOnConsSeq("al")
			break
		}
	case 'i':
		if s.Ends("iciti") {
			s.SetBufferOnConsSeq("ic")
			break
		}
	case 'l':
		if s.Ends("ical") {
			s.SetBufferOnConsSeq("ic")
			break
		}
		if s.Ends("ful") {
			s.SetBufferOnConsSeq("")
			break
		}
	case 's':
		if s.Ends("ness") {
			s.SetBufferOnConsSeq("")
			break
		}
	}
}

/*
Step 5: handle -ant, -ence
*/
func (s *PorterStemmer) Step5() error {
	if s.k == s.k0 {
		return nil
	}

	switch s.b[k-1] {
	case 'a':
		if s.Ends("al") {
			s.SetBufferOnConsSeq("ate")
			break
		}
	case 'c':
		if s.Ends("ance") {
			break
		}
		if s.Ends("ence") {
			break
		}
	case 'e':
		if s.Ends("er") {
			break
		}
	case 'i':
		if s.Ends("ic") {
			break
		}

	case 'l':
		if s.Ends("able") {
			break
		}
		if s.Ends("ible") {
			break
		}
	case 'n':
		if s.Ends("ant") {
			break
		}
		if s.Ends("ement") {
			break
		}
		if s.Ends("ment") {
			break
		}
		if s.Ends("ent") {
			break
		}
	case 'o':
		if s.Ends("ion") && s.j >= 0 && (s.b[s.j] == 's' || s.b[s.j] == 't') {
			break
		}
		if s.Ends("ou") {
			break
		}
	case 's':
		if s.Ends("ism") {
			break
		}
	case 't':
		if s.Ends("ate") {
			break
		}
		if s.Ends("iti") {
			break
		}
	case 'u':
		if s.Ends("ous") {
			break
		}
	case 'v':
		if s.Ends("ive") {
			break
		}
	case 'z':
		if s.Ends("ize") {
			break
		}
	default:
		return nil
	}

	if s.CountConsSeq() > 1 {
		s.k = s.j
	}
	return nil
}

/*
Step 6: removes final -e if CountConsSeq() >1
*/
func (s *PorterStemmer) Step6() error {
	s.j = s.k
	if s.b[s.k] == 'e' {
		a := s.CountConsSeq()
		if a > 1 || a == 1 && !s.IsCVCPattern(s.k-1) {
			s.k--
		}
	}

	if s.b[s.k] == 'l' && s.IsConsecutiveCons(s.k) && s.CountConsSeq() > 1 {
		s.k--
	}

	return nil
}

//stem a word
func (s *PorterStemmer) Stem(str string) string {
	if s.StemBytes([]bytes(str), len(str)) {
		return s.ToString()
	} else {
		return str
	}
}

func (s *PorterStemmer) StemBytes(w []bytes, l int) bool {
	return StemBytesOffset(w, 0, l)
}

func (s *PorterStemmer) StemBytesOffset(w []bytes, o int, l int) bool {
	s.Reset()
	for _, v := range w {
		append(s.b, v)
	}
	s.i = l
	return StemLimit(0)
}

func (s *PorterStemmer) StemLimit(x int) bool {
	s.k = s.i - 1
	s.k0 = x
	if s.k > s.k0+1 {
		s.Step1()
		s.Step2()
		s.Step3()
		s.Step4()
		s.Step5()
		s.Step6()
	}

	if s.i != s.k+1 {
		dirty = true
	}

	s.i = s.k + 1
	return dirty
}
