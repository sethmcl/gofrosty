package nsemver

import (
	"testing"
)

// Define test data for ParseVersion()
type parseVersionTest struct {
	input    string
	expected string
}

var parseVersionTests = []parseVersionTest{
	parseVersionTest{"=1.100.30", "1.100.30"},
	parseVersionTest{"=v0.0.0-preA", "0.0.0-preA"},
	parseVersionTest{"=v0.0.0-preA.5", "0.0.0-preA.5"},
	parseVersionTest{"=v0.0.0.1.2-preA", "ERROR"},
	parseVersionTest{"v0.0.0", "0.0.0"},
	parseVersionTest{"1.x.X", "ERROR"},
	parseVersionTest{"~1.0.0", "ERROR"},
	parseVersionTest{"^1.0.0", "ERROR"},
	parseVersionTest{"1.0.0 || 2.0.0", "ERROR"},
	parseVersionTest{"4.0", "ERROR"},
	parseVersionTest{">1.0.0", "ERROR"},
	parseVersionTest{"4", "ERROR"},
}

func TestParseVersion(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range parseVersionTests {
		actual, err := ParseVersion(data.input)
		if data.expected == "ERROR" {
			test.Assert(err != nil, true, err.Error())
		} else {
			test.Assert(err, nil)
			test.Assert(actual.String(), data.expected)
		}
	}
}

// Define test data for Version.Compare()
type compareTest struct {
	a        string
	b        string
	expected int
}

var compareTests = []compareTest{
	compareTest{"1.0.0", "1.0.0", 0},
	compareTest{"1.0.1", "1.0.0", 1},
	compareTest{"1.1.0", "1.0.0", 1},
	compareTest{"2.0.0", "1.0.0", 1},
	compareTest{"1.0.0", "1.0.100", -1},
	compareTest{"1.0.0", "1.1.0", -1},
	compareTest{"1.0.0", "2.0.0", -1},
	compareTest{"1.0.0", "1.0.0-alpha", 1},
	compareTest{"1.0.0-alpha", "1.0.0", -1},
	compareTest{"1.0.0-alpha5", "1.0.0-alpha5", 0},
	compareTest{"1.0.0-alpha6", "1.0.0-alpha5", 1},
	compareTest{"1.0.0-alpha5", "1.0.0-alpha6", -1},
	compareTest{"1.0.0-alpha5", "1.0.0-beta5", -1},
	compareTest{"1.0.0-beta5", "1.0.0-alpha5", 1},
}

func TestCompareVersions(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range compareTests {
		va, err := ParseVersion(data.a)
		test.Assert(err, nil)
		vb, err := ParseVersion(data.b)
		test.Assert(err, nil)
		test.Assert(va.Compare(vb), data.expected)
	}
}

// Define test data for Version.{Eq,Lt,Gt,Gte,Lte}()
type compareHelperTest struct {
	a        string
	b        string
	operator string
	expected bool
}

var compareHelperTests = []compareHelperTest{
	compareHelperTest{"1.0.0", "1.0.0", "Eq", true},
	compareHelperTest{"1.0.0", "1.0.0", "Lte", true},
	compareHelperTest{"1.0.0", "1.0.0", "Gte", true},
	compareHelperTest{"1.0.0", "1.0.0-alpha", "Eq", false},
	compareHelperTest{"1.0.1", "1.0.0", "Lt", false},
	compareHelperTest{"1.0.1", "1.0.0", "Lte", false},
	compareHelperTest{"1.0.1", "1.0.0", "Gt", true},
	compareHelperTest{"1.0.1", "1.0.0", "Gte", true},
}

func TestCompareHelpers(t *testing.T) {
	test := NewTestUtil(t)

	for _, data := range compareHelperTests {
		va, err := ParseVersion(data.a)
		test.Assert(err, nil)
		vb, err := ParseVersion(data.b)
		test.Assert(err, nil)
		switch data.operator {
		case "Eq":
			test.Assert(va.Eq(vb), data.expected, data.a, data.operator, data.b)
		case "Lt":
			test.Assert(va.Lt(vb), data.expected, data.a, data.operator, data.b)
		case "Gt":
			test.Assert(va.Gt(vb), data.expected, data.a, data.operator, data.b)
		case "Lte":
			test.Assert(va.Lte(vb), data.expected, data.a, data.operator, data.b)
		case "Gte":
			test.Assert(va.Gte(vb), data.expected, data.a, data.operator, data.b)
		default:
			t.Errorf("Invalid operator specified in test: %s", data.operator)
		}
	}
}
