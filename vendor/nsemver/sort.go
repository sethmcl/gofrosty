package nsemver

import (
	"strings"
)

// Versions group of *Versions
type Versions []*Version

// NewVersions create a new Versions from []string
func NewVersions(vs []string) (Versions, error) {
	var result Versions
	for _, v := range vs {
		ver, err := ParseVersion(v)
		if err != nil {
			return nil, err
		}
		result = append(result, ver)
	}
	return result, nil
}

// Len returns length of Versions slice
func (slice Versions) Len() int {
	return len(slice)
}

// Less returns comparasion
func (slice Versions) Less(a, b int) bool {
	return slice[a].Gt(slice[b])
}

// Swap swaps the position of two items
func (slice Versions) Swap(a, b int) {
	slice[a], slice[b] = slice[b], slice[a]
}

// String returns string representation
func (slice Versions) String() string {
	var strs []string
	for _, v := range slice {
		strs = append(strs, v.String())
	}
	return strings.Join(strs, ", ")
}
