//package con holds all the types used in the library
package util

import (
	"strconv"
	"strings"
)

//structure for storing start and end integer offset
type Span struct {
	start, end int64
	class      string
}

//allocates and return a new Span
func NewSpan(st int64, ed int64) *Span {
	if st < 0 {
		panic("start index should be >= 0")
	}
	if ed < 0 {
		panic("end index must be >= 0")
	}
	if st > ed {
		panic("start index must not be larger than end index")
	}
	s := Span{st, ed, ""}

	return &s
}

//specify type to the Span
func (s *Span) AddClass(c string) *Span {
	s.class = c

	return s
}

//add offset to the Span
func (s *Span) AddOffset(offset int64) *Span {
	s.start = s.start + offset
	s.end = s.end + offset

	return s
}

//get the start of the Span
func (s *Span) GetStart() int64 {
	return s.start
}

//get the end of the Span
func (s *Span) GetEnd() int64 {
	return s.end
}

//get the class of the Span
func (s *Span) GetClass() string {
	return s.class
}

//get the length of the Span
func (s *Span) GetLength() int64 {
	return s.end - s.start
}

//check if this Span contains the specified Span
func (s *Span) ContainsSpan(sp Span) bool {
	return s.start <= sp.GetStart() && sp.GetEnd() <= s.end
}

//check if this Span contains the specified index
//if the index value is equal to end of this Span then it is
//considered outside
func (s *Span) ContainsIndex(i int64) bool {
	return s.start <= i && i < s.end
}

//check if this Span starts with the specified Span as
//well as contains the Span
func (s *Span) StartsWithSpan(sp Span) bool {
	return s.start == sp.GetStart() && s.ContainsSpan(sp)
}

//check if this Span intersect with the specified Span
func (s *Span) Intersects(sp Span) bool {
	return s.ContainsSpan(sp) || sp.ContainsSpan(*s) ||
		s.GetStart() <= sp.GetStart() && sp.GetStart() < s.GetEnd() ||
		sp.GetStart() <= s.GetStart() && s.GetStart() < sp.GetEnd()
}

//check if this Span crosses the specified Span
func (s *Span) Crosses(sp Span) bool {
	return !s.ContainsSpan(sp) || !sp.ContainsSpan(*s) ||
		s.GetStart() <= sp.GetStart() && sp.GetStart() < s.GetEnd() ||
		sp.GetStart() <= s.GetStart() && s.GetStart() < sp.GetEnd()
}

//get the string covered by this Span
func (s *Span) GetCoveredText(str string) string {
	strByte := []byte(str)
	if s.GetEnd() > int64(len(strByte)) {
		panic("Span not valid for given string")
	}

	return string(strByte[s.GetStart():s.GetEnd()])
}

//compare this Span with the specified Span
func (s *Span) CompareSpan(sp Span) int64 {
	if s.GetStart() < sp.GetStart() {
		return -1
	} else if s.GetStart() == sp.GetStart() {
		if s.GetEnd() > sp.GetEnd() {
			return -1
		} else if s.GetEnd() < sp.GetEnd() {
			return 1
		} else {
			// compare the type
			if s.GetClass() == "" && sp.GetClass() == "" {
				return 0
			} else if s.GetClass() != "" && sp.GetClass() != "" {
				// use type lexicography order
				//TODO: handle getType().compareTo(s.getType())
				if s.GetClass() == sp.GetClass() {
					return 0
				} else {
					return -1
				}
			} else if s.GetClass() != "" {
				return -1
			}
			return 1
		}
	} else {
		return 1
	}
}

//get the string form of the Span
func (s *Span) ToString() string {
	var str []string

	str = append(str, "[")
	str = append(str, strconv.Itoa(int(s.GetStart())))
	str = append(str, "...")
	str = append(str, strconv.Itoa(int(s.GetEnd())))
	str = append(str, ")")
	if s.GetClass() != "" {
		str = append(str, " ")
		str = append(str, s.GetClass())
	}

	return strings.Join(str, "")
}

//convert array of spans to array of strings
func SpansToStrings(spans []*Span, str string) (tokens []string) {
	for _, v := range spans {
		tokens = append(tokens, v.GetCoveredText(str))
	}
	return
}
