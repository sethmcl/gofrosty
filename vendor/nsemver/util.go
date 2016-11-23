package nsemver

import (
	"fmt"
	"sort"
	"strings"
)

// Clean removes extraneous characters from semver string
//
//     "v0.0.0" => "0.0.0"
//     "V0.0.0" => "0.0.0"
//     "0. 0.0" => "0.0.0"
func Clean(v string) string {
	out := ""
	for _, c := range v {
		if c != 'v' && c != 'V' && c != ' ' {
			out += string(c)
		}
	}
	return out
}

// ComparePR compares two Prerelease (PR) version values
// Detailed description of comparison logic: https://docs.npmjs.com/misc/semver#prerelease-tags
func ComparePR(a string, b string) int {
	if a == b {
		return 0
	}

	if a == "" && b != "" {
		return 1
	}

	if a != "" && b == "" {
		return -1
	}

	return strings.Compare(a, b)
}

// MatchLatest return latest version that satisfies provided range
func MatchLatest(input string, candidates []string) (string, error) {
	vs, err := NewVersions(candidates)
	if err != nil {
		return "", err
	}

	// Sort (default sort is newest -> oldest)
	sort.Sort(vs)

	rang, err := ParseRange(input)
	if err != nil {
		return "", err
	}

	for _, ver := range vs {
		if rang.SatisfiedBy(ver) {
			return ver.Raw, nil
		}
	}

	return "", fmt.Errorf("%s not satisfied by any candidates: %s", input, candidates)
}

// PadVersion adds missing parts to version string and replaces x|X|* with 0
func PadVersion(v string) string {
	pv := ""
	for _, r := range v {
		if isXRune(r) {
			pv += "0"
		} else {
			pv += string(r)
		}

	}

	len := len(strings.Split(pv, "."))

	if len == 1 {
		pv += ".0.0"
	}

	if len == 2 {
		pv += ".0"
	}

	return pv
}
