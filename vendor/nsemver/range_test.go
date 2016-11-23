package nsemver

import (
	"testing"
)

// Define test data for ParseRange()
type parseRangeTest struct {
	input    string
	expected string
}

var parseRangeTests = []parseRangeTest{
	parseRangeTest{"1.0.x", ">=1.0.0 <1.1.0"},
	parseRangeTest{"2.0.0", "=2.0.0"},
	parseRangeTest{"~1.1", ">=1.1.0 <1.2.0"},
	parseRangeTest{"^ 1.1", ">=1.1.0 <2.0.0"},
	parseRangeTest{"^ 1.1.0", ">=1.1.0 <2.0.0"},
	parseRangeTest{"~ 1.1", ">=1.1.0 <1.2.0"},
	parseRangeTest{"~ 0.1.11", ">=0.1.11 <0.2.0"},
	parseRangeTest{"^3.1.0", ">=3.1.0 <4.0.0"},
	parseRangeTest{"~0.10.x", ">=0.10.0 <0.11.0"},
	parseRangeTest{"^0.10.x", ">=0.10.0 <0.11.0"},
	parseRangeTest{"~2.0.0", ">=2.0.0 <2.1.0"},
	parseRangeTest{"^2.0.0", ">=2.0.0 <3.0.0"},
	parseRangeTest{"2", ">=2.0.0 <3.0.0"},
	parseRangeTest{"3.1", ">=3.1.0 <3.2.0"},
	parseRangeTest{"2||3", ">=2.0.0 <3.0.0||>=3.0.0 <4.0.0"},
	parseRangeTest{"2 || 3", ">=2.0.0 <3.0.0||>=3.0.0 <4.0.0"},
	parseRangeTest{"2.0 || 3", ">=2.0.0 <2.1.0||>=3.0.0 <4.0.0"},
	parseRangeTest{">=1.0.0 <2.0.0 || >3.0.0 <6.0.0", ">=1.0.0 <2.0.0||>3.0.0 <6.0.0"},
	parseRangeTest{">=2.6.1 <4", ">=2.6.1 <4.0.0"},
	parseRangeTest{"*", ">=0.0.0"},
	parseRangeTest{"latest", ">=0.0.0"},
	parseRangeTest{">= 1.5.x", ">=1.5.0"},
}

func TestParseRange(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range parseRangeTests {
		actual, err := ParseRange(data.input)
		if data.expected == "ERROR" {
			test.Assert(err != nil, true)
		} else {
			test.Assert(err, nil)
			if err == nil {
				test.Assert(actual.String(), data.expected)
			}
		}
	}
}

// Define test data for Parse*RangeSet()
type parseRangeSetTest struct {
	input    string
	expected string
	rtype    string
}

var parseRangeSetTests = []parseRangeSetTest{
	parseRangeSetTest{"1.0.x", ">=1.0.0 <1.1.0", "x"},
	parseRangeSetTest{"1.x.x", ">=1.0.0 <2.0.0", "x"},
	parseRangeSetTest{"1.X.x", ">=1.0.0 <2.0.0", "x"},
	parseRangeSetTest{"1.*.x", ">=1.0.0 <2.0.0", "x"},
	parseRangeSetTest{"1.x", ">=1.0.0 <2.0.0", "x"},
	parseRangeSetTest{"*", ">=0.0.0", "x"},
	parseRangeSetTest{"x", ">=0.0.0", "x"},
	parseRangeSetTest{"1.0.x", "error", "~"},
	parseRangeSetTest{"1.0.x", "error", "^"},
}

func TestParseRangeSet(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range parseRangeSetTests {
		var (
			actual Set
			err    error
		)

		switch data.rtype {
		case "x":
			actual, err = ParseXRangeSet(data.input)
		case "^":
			actual, err = ParseCaretRangeSet(data.input)
		case "~":
			actual, err = ParseTildeRangeSet(data.input)
		}

		if data.expected == "error" {
			test.Assert(err != nil, true, data.input)
		} else {
			test.Assert(err, nil)
			test.Assert(actual.String(), data.expected, data.input)
		}
	}
}

// Define test data for SatisfiedBy()
type satisfiedByTest struct {
	rang     string
	ver      string
	expected bool
}

var satisfiedByTests = []satisfiedByTest{
	satisfiedByTest{">=1.0.0 <2.0.0", "1.5.0", true},
	satisfiedByTest{">=1.0.0 <2.0.0 || >3.0.0 <6.0.0", "1.5.0", true},
	satisfiedByTest{">=1.0.0 <2.0.0 || >3.0.0 <6.0.0", "2.5.0", false},
	satisfiedByTest{">=1.0.0 <2.0.0 || >3.0.0 <6.0.0", "3.5.0", true},
	satisfiedByTest{"2.x||4.1.x", "3.5.0", false},
	satisfiedByTest{"2.x||4.1.x", "2.9.0", true},
	satisfiedByTest{"2.x||4.1.x", "4.1.0", true},
	satisfiedByTest{"2.x||4.1.x", "4.1.10", true},
	satisfiedByTest{"2.x||4.1.x", "4.2.0", false},
	satisfiedByTest{"~2.0.0||4.1.x", "4.2.0", false},
	satisfiedByTest{"2 || 3 || 4", "4.2.0", true},
	satisfiedByTest{"2 || 5 || 4", "4.2.0", true},
	satisfiedByTest{"2 || 5 || 3", "4.2.0", false},
}

func TestSatisfiedBy(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range satisfiedByTests {
		var (
			actual bool
			rang   *Range
			ver    *Version
			err    error
		)

		rang, err = ParseRange(data.rang)
		test.Assert(err, nil)

		ver, err = ParseVersion(data.ver)
		test.Assert(err, nil)

		actual = rang.SatisfiedBy(ver)
		test.Assert(actual, data.expected, data.rang, data.ver)
	}
}
