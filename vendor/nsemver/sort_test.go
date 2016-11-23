package nsemver

import (
	"sort"
	"testing"
)

// Define test data for Sort()
type sortTest struct {
	input    []string
	expected string
}

var sortTests = []sortTest{
	sortTest{[]string{"0.2.0", "3.2.1", "3.6.1", "3.1.0"}, "3.6.1, 3.2.1, 3.1.0, 0.2.0"},
}

func TestSort(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range sortTests {
		vs, err := NewVersions(data.input)
		test.Assert(err, nil)
		sort.Sort(vs)
		test.Assert(vs.String(), data.expected)
	}
}
