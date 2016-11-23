package nsemver

import (
	"strconv"
)

// RuneBuffer builds up a string from runes
type RuneBuffer struct {
	data string
}

// NewRuneBuffer creates a new *RuneBuffer
func NewRuneBuffer() *RuneBuffer {
	return &RuneBuffer{""}
}

// Append adds a rune to the buffer
func (b *RuneBuffer) Append(r rune) {
	b.data += string(r)
}

// String converts buffer to a string
func (b *RuneBuffer) String() string {
	return b.data
}

// Int converts buffer to an int
func (b *RuneBuffer) Int() (int, error) {
	return strconv.Atoi(b.data)
}

// Clear resets the buffer
func (b *RuneBuffer) Clear() {
	b.data = ""
}

// Len returns length of buffer (number of runes)
func (b *RuneBuffer) Len() int {
	return len(b.data)
}
