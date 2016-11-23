package nsemver

import (
	"fmt"
)

// Version represents a parsed semantic version string
type Version struct {
	Raw   string
	Major int
	Minor int
	Patch int
	PR    string
}

func (ver *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
	if len(ver.PR) > 0 {
		return fmt.Sprintf("%s-%s", s, ver.PR)
	}
	return s
}

// Compare two *Version references. Returns -1, 0, or 1.
func (ver *Version) Compare(b *Version) int {
	if ver.Major < b.Major {
		return -1
	}

	if ver.Major > b.Major {
		return 1
	}

	if ver.Minor < b.Minor {
		return -1
	}

	if ver.Minor > b.Minor {
		return 1
	}

	if ver.Patch < b.Patch {
		return -1
	}

	if ver.Patch > b.Patch {
		return 1
	}

	return ComparePR(ver.PR, b.PR)
}

// Eq returns true if the versions are equal
func (ver *Version) Eq(b *Version) bool {
	return ver.Compare(b) == 0
}

// Lt returns true if ver is less than b
func (ver *Version) Lt(b *Version) bool {
	return ver.Compare(b) == -1
}

// Gt returns true if ver is greater than than b
func (ver *Version) Gt(b *Version) bool {
	return ver.Compare(b) == 1
}

// Lte returns true if ver is less than or equal to b
func (ver *Version) Lte(b *Version) bool {
	c := ver.Compare(b)
	return (c == -1) || (c == 0)
}

// Gte returns true if ver is greater than or equal to b
func (ver *Version) Gte(b *Version) bool {
	c := ver.Compare(b)
	return (c == 1) || (c == 0)
}

// ParseVersion parses a string and return a Version object
func ParseVersion(v string) (*Version, error) {
	ver := &Version{Raw: v}

	const (
		stateParseMajor = 1
		stateParseMinor = 2
		stateParsePatch = 3
		stateParsePR    = 5
	)

	missingMajorErr := fmt.Errorf("%s is missing major version", v)
	missingMinorErr := fmt.Errorf("%s is missing minor version", v)
	missingPatchErr := fmt.Errorf("%s is missing patch version", v)

	state := stateParseMajor
	buffer := NewRuneBuffer()

	// examine each rune in the version string
	err := IterateString(v, func(r rune, iter *StringIteratorState) error {
		// if current rune indicates a range value, then bail
		if isRangeRune(r) {
			return fmt.Errorf("%s is not an explicit version [found %s]", v, string(r))
		}

		switch state {
		case stateParseMajor:
			if isDigitRune(r) {
				buffer.Append(r)
			}

			next, _ := iter.Next()
			if iter.HasNext() && isVersionSeparatorRune(next) {
				// We have reached end of major version. Let's try to convert
				// it to an int and record it.
				if buffer.Len() == 0 {
					return missingMajorErr
				}

				i, err := buffer.Int()
				if err != nil {
					return err
				}
				ver.Major = i
				buffer.Clear()
				state = stateParseMinor
			}

			if !iter.HasNext() {
				// We have reached end of string prematurely
				if buffer.Len() == 0 {
					return missingMajorErr
				}
				return missingMinorErr
			}

			return nil
		case stateParseMinor:
			if isDigitRune(r) {
				buffer.Append(r)
			}

			next, _ := iter.Next()
			if iter.HasNext() && isVersionSeparatorRune(next) {
				// We have reached end of minor version. Let's try_url
				// to convert it to an int and record it.
				if buffer.Len() == 0 {
					return missingMinorErr
				}

				i, err := buffer.Int()
				if err != nil {
					return err
				}
				ver.Minor = i
				buffer.Clear()
				state = stateParsePatch
			}

			if !iter.HasNext() {
				// We have reached end of string prematurely
				if buffer.Len() == 0 {
					return missingMinorErr
				}
				return missingPatchErr
			}

			return nil
		case stateParsePatch:
			if isDigitRune(r) {
				buffer.Append(r)
			}

			next, _ := iter.Next()
			if iter.HasNext() && next == '.' {
				return fmt.Errorf("%s is invalid. Must be in form of MAJOR.MINOR.PATCH-PR", v)
			}

			if !iter.HasNext() || next == '-' {
				// We have reached end of minor version. Let's try_url
				// to convert it to an int and record it.
				if buffer.Len() == 0 {
					return missingPatchErr
				}

				i, err := buffer.Int()
				if err != nil {
					return err
				}
				ver.Patch = i
				buffer.Clear()

				// skip over first '-'
				iter.Skip(1)
				state = stateParsePR
			}

			return nil
		case stateParsePR:
			buffer.Append(r)
			if !iter.HasNext() {
				ver.PR = buffer.String()
			}

			return nil
		}

		return nil
	})

	if err != nil {
		return ver, err
	}

	return ver, nil
}
