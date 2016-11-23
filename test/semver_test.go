package test

import (
	"github.com/sethmcl/gofrosty/lib"
	"github.com/sethmcl/gofrosty/lib/testutil"
	"github.com/sethmcl/gofrosty/vendor/semver"
	"testing"
)

// TestConvertNpmSemver
func TestConvertNpmSemver(t *testing.T) {
	test := testutil.New(t)

	table := map[string]string{
		// "1.0.0":       "1.0.0",
		// "1.0.x":       ">=1.0.0 <1.1.0",
		// "1.0.X":       ">=1.0.0 <1.1.0",
		// "v1.0.X":      ">=1.0.0 <1.1.0",
		// "~1.0.0":      ">=1.0.0 <1.1.0",
		// "1.3.X":       ">=1.3.0 <1.4.0",
		// "1.0":         ">=1.0.0 <1.1.0",
		// "2 || 3 || 4": ">=2.0.0 <3.0.0 || >=3.0.0 <4.0.0 || >=4.0.0 <5.0.0",
		// "^1.0.27-1":   ">=1.0.27 <2.0.0",
		// "^1.5.27-1":   ">=1.5.27 <2.0.0",
		// "~1.5.27-1":   ">=1.5.27 <1.6.0",
		">=2.6.1 <4": ">=2.6.1 <4.0.0",
	}

	for input, expected := range table {
		actual, err := lib.ConvertNpmSemver(input)
		test.Assert(err, nil)
		test.Assert(actual, expected)

		_, err = semver.ParseRange(actual)
		test.Assert(err, nil)
	}
}

// // TestIsExplicitSemver
// func TestIsExplicitSemver(t *testing.T) {
// 	test := testutil.New(t)

// 	table := map[string]bool{
// 		"1.0.0":        true,
// 		"v0.100.0":     true,
// 		"1.0.x":        false,
// 		"1.0.X":        false,
// 		"v1.0.X":       false,
// 		"~1.0.0":       false,
// 		"1.3.X":        false,
// 		"1.0":          false,
// 		"1.1.1.1":      false,
// 		">= 1.3.0 < 2": false,
// 	}

// 	for input, expected := range table {
// 		actual := lib.IsExplicitSemver(input)
// 		test.Assert(actual, expected)
// 	}
// }

// // TestIsSemverRange
// func TestIsSemverRange(t *testing.T) {
// 	test := testutil.New(t)

// 	table := map[string]bool{
// 		"1.0.0":                            false,
// 		"v0.100.0":                         false,
// 		"1.0.x":                            true,
// 		"1.0.X":                            true,
// 		"v1.0.X":                           true,
// 		"~1.0.0":                           true,
// 		"1.3.X":                            true,
// 		"1.0":                              true,
// 		"1.1.1.1":                          false,
// 		"=1.0.0":                           true,
// 		">1.0.0 <2.0.0":                    true,
// 		"2 || 3 || 4":                      true,
// 		">=2.0.0 <3.0.0 || >=4.0.0 <5.0.0": true,
// 	}

// 	for input, expected := range table {
// 		actual := lib.IsSemverRange(input)
// 		test.Assert(actual, expected, input)
// 	}
// }

// // // TestSanitizeSemver
// // func TestSanitizeSemver(t *testing.T) {
// // 	test := testutil.New(t)

// // 	table := map[string]string{
// // 		"4":             "4.0.0",
// // 		"4.0":           "4.0.0",
// // 		"4.0.0":         "4.0.0",
// // 		"4.0.0.0":       "4.0.0.0",
// // 		">=3.0 <4":      ">=3.0.0 <4.0.0",
// // 		">= 3.0 < 4":    ">=3.0.0 <4.0.0",
// // 		">= 3.0 <    4": ">=3.0.0 <4.0.0",
// // 	}

// // 	for input, expected := range table {
// // 		actual := lib.SanitizeSemver(input)
// // 		test.Assert(actual, expected, input)
// // 	}
// // }

// // TestResolveSemver
// func TestResolveSemver(t *testing.T) {
// 	test := testutil.New(t)

// 	type TestTableEntry struct {
// 		Semver     string
// 		Candidates []string
// 		Expected   string
// 	}

// 	table := []TestTableEntry{
// 		TestTableEntry{
// 			Semver:     "1.0.x",
// 			Candidates: []string{"1.0.0", "0.9.0-pre"},
// 			Expected:   "1.0.0",
// 		},
// 		TestTableEntry{
// 			Semver:     "!1.!!0.x",
// 			Candidates: []string{"1.0.0", "0.9.0-pre"},
// 			Expected:   "error",
// 		},
// 		TestTableEntry{
// 			Semver:     "~2.5.0",
// 			Candidates: []string{"1.0.0", "2.5.0", "2.5.5", "0.9.0-pre"},
// 			Expected:   "2.5.5",
// 		},
// 		TestTableEntry{
// 			Semver:     "~2.5.0",
// 			Candidates: []string{"1.0.0", "2.5.5", "2.5.0", "0.9.0-pre"},
// 			Expected:   "2.5.5",
// 		},
// 		TestTableEntry{
// 			Semver: "1.x.x",
// 			Candidates: []string{
// 				"1.1.0",
// 				"2.0.0",
// 				"2.0.1",
// 				"2.0.2",
// 				"1.0.0",
// 				"1.0.1",
// 				"1.0.2",
// 				"1.0.3"},
// 			Expected: "1.1.0",
// 		},
// 	}

// 	for _, entry := range table {
// 		v, err := lib.ResolveSemver(entry.Semver, entry.Candidates)
// 		if entry.Expected == "error" {
// 			isError := (err != nil)
// 			test.Assert(isError, true)
// 		} else {
// 			test.Assert(err, nil)
// 			test.Assert(v, entry.Expected)
// 		}
// 	}
// }
