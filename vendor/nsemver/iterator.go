package nsemver

import (
	"fmt"
)

// StringIterator iterates over a string, rune by rune
type StringIterator struct {
	Target  string
	Index   int
	Min     int
	Max     int
	Visitor StringVisitor
}

// IterateString create a new StringIterator
func IterateString(s string, f StringVisitor) error {
	si := &StringIterator{
		Target:  s,
		Index:   0,
		Min:     0,
		Max:     len(s),
		Visitor: f,
	}

	return si.Run()
}

// Run the iterator
func (si *StringIterator) Run() error {
	state := &StringIteratorState{iterator: si}
	for si.Index = 0; si.Index < si.Max; si.Index++ {
		err := si.Visitor(rune(si.Target[si.Index]), state)
		if err != nil {
			return err
		}
		if state.requestAbort {
			break
		}
	}
	return nil
}

// StringVisitor function to visit a rune, called by StringIterator
type StringVisitor func(rune, *StringIteratorState) error

// StringIteratorState exposes iterator state to the visitor function
type StringIteratorState struct {
	iterator     *StringIterator
	requestAbort bool
}

// HasNext returns true if we have not reached end of string
func (sis *StringIteratorState) HasNext() bool {
	return sis.iterator.Index < (sis.iterator.Max - 1)
}

// Next returns next rune in string
func (sis *StringIteratorState) Next() (rune, error) {
	if sis.HasNext() {
		return rune(sis.iterator.Target[sis.iterator.Index+1]), nil
	}

	return 0, fmt.Errorf("End of string error")
}

// Abort halts the iterator
func (sis *StringIteratorState) Abort() {
	sis.requestAbort = true
}

// HasPrev returns true if we are not at very start of string
func (sis *StringIteratorState) HasPrev() bool {
	return sis.iterator.Index > sis.iterator.Min
}

// Skip skips over next rune by incrementing index value
func (sis *StringIteratorState) Skip(i int) {
	sis.iterator.Index += i
}
