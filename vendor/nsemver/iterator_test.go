package nsemver

import (
	"fmt"
	"testing"
)

func TestStringIterator(t *testing.T) {
	test := NewTestUtil(t)

	visited := ""
	target := "1234567890"
	err := IterateString(target, func(r rune, iter *StringIteratorState) error {
		visited += string(r)
		if !iter.HasNext() {
			visited += "--end!--"
		}
		return nil
	})
	test.Assert(err, nil)
	test.Assert(visited, "1234567890--end!--")

	visited = ""
	err = IterateString(target, func(r rune, iter *StringIteratorState) error {
		if !iter.HasPrev() {
			visited += "(start)"
		}
		visited += string(r)
		return nil
	})
	test.Assert(err, nil)
	test.Assert(visited, "(start)1234567890")

	visited = ""
	err = IterateString(target, func(r rune, iter *StringIteratorState) error {
		visited += string(r)
		if r == '5' {
			iter.Abort()
		}
		return nil
	})
	test.Assert(err, nil)
	test.Assert(visited, "12345")

	visited = ""
	err = IterateString(target, func(r rune, iter *StringIteratorState) error {
		visited += string(r)
		if r == '5' {
			return fmt.Errorf("why so blue, %s?", string(r))
		}
		return nil
	})
	test.Assert(visited, "12345")
	test.Assert(err.Error(), "why so blue, 5?")
}
