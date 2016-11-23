package nsemver

import (
	"testing"
)

// Define test data for CleanSemver()
type cleanSemverTest struct {
	input    string
	expected string
}

var cleanTests = []cleanSemverTest{
	cleanSemverTest{"v0.0.0", "0.0.0"},
	cleanSemverTest{"V0.0.0", "0.0.0"},
	cleanSemverTest{"0. 0.0", "0.0.0"},
	cleanSemverTest{">=2.0.0 <3.0.0", ">=2.0.0<3.0.0"},
}

func TestClean(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range cleanTests {
		actual := Clean(data.input)
		test.Assert(actual, data.expected)
	}
}

// Define test data for MatchLatest()
type matchLatestTest struct {
	input      string
	candidates []string
	expected   string
}

var matchLatestTests = []matchLatestTest{
	matchLatestTest{">=1.0.0", []string{"0.9.0", "1.2.0", "1.1.0"}, "1.2.0"},
	matchLatestTest{"1.1.0", []string{"0.9.0", "1.2.0", "1.1.0"}, "1.1.0"},
	matchLatestTest{"1.x", []string{"0.9.0", "1.2.0", "1.1.0", "3.0.0"}, "1.2.0"},
	matchLatestTest{"1.x-alpha", []string{"0.9.0", "1.2.0", "1.2.0-alpha", "1.1.0", "3.0.0"}, "1.2.0"},
	matchLatestTest{"^3.1.0", []string{"0.2.0", "3.2.1", "3.6.1", "3.1.0"}, "3.6.1"},
}

func TestMatchLatest(t *testing.T) {
	test := NewTestUtil(t)
	for _, data := range matchLatestTests {
		actual, err := MatchLatest(data.input, data.candidates)
		if data.expected == "error" {
			test.Assert(err != nil, true, err.Error())
		} else {
			test.Assert(err, nil)
			test.Assert(actual, data.expected, data.input)
		}
	}
}

// Define test data for PadVersion()
type padVersionTest struct {
	input    string
	expected string
}

var padVersionTests = []padVersionTest{
	padVersionTest{"1", "1.0.0"},
	padVersionTest{"1.0", "1.0.0"},
	padVersionTest{"1.0.0", "1.0.0"},
	padVersionTest{"1.0.0.0", "1.0.0.0"},
	padVersionTest{"1.*.0.0", "1.0.0.0"},
	padVersionTest{"1.0.*", "1.0.0"},
	padVersionTest{"1.0.x", "1.0.0"},
	padVersionTest{"1.0.X", "1.0.0"},
	padVersionTest{"1.*", "1.0.0"},
	padVersionTest{"1.x", "1.0.0"},
	padVersionTest{"1.X", "1.0.0"},
	padVersionTest{"*", "0.0.0"},
	padVersionTest{"x", "0.0.0"},
	padVersionTest{"X", "0.0.0"},
}

func TestPadVersion(t *testing.T) {
	test := NewTestUtil(t)
	for _, data := range padVersionTests {
		test.Assert(PadVersion(data.input), data.expected, data.input)
	}
}
