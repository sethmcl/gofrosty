package nsemver

import (
	"fmt"
	"strconv"
	"strings"
)

// Range represents a parsed semantic version range
type Range struct {
	Raw  string
	Sets []Set
}

// Set is a group of constraints
type Set []*Constraint

func (set Set) String() string {
	s := ""
	for cidx, constraint := range set {
		s += constraint.String()
		if cidx < len(set)-1 {
			s += " "
		}
	}
	return s
}

// SatisfiedBy returns true if the version satisfies this set of constraints
func (set Set) SatisfiedBy(ver *Version) bool {
	for _, constraint := range set {
		switch constraint.Operator {
		case "=":
			pass := ver.Eq(constraint.Version)
			if !pass {
				return false
			}
		case ">":
			pass := ver.Gt(constraint.Version)
			if !pass {
				return false
			}
		case ">=":
			pass := ver.Gte(constraint.Version)
			if !pass {
				return false
			}
		case "<":
			pass := ver.Lt(constraint.Version)
			if !pass {
				return false
			}
		case "<=":
			pass := ver.Lte(constraint.Version)
			if !pass {
				return false
			}
		}
	}

	return true
}

// AddSet to range
func (r *Range) AddSet(s Set) {
	r.Sets = append(r.Sets, s)
}

// String returns string representation of range
func (r *Range) String() string {
	s := ""
	for sidx, set := range r.Sets {
		if len(set) == 0 {
			continue
		}

		s += set.String()
		if sidx < len(r.Sets)-1 {
			s += "||"
		}
	}
	return s
}

// SatisfiedBy returns true if this version satisfies the range constraints
// a range has multiple sets of constraints
// result will be true if version satisfies any of these sets
func (r *Range) SatisfiedBy(ver *Version) bool {
	for _, set := range r.Sets {
		if set.SatisfiedBy(ver) {
			return true
		}
	}

	return false
}

// Constraint represents a single condition within a range expression
//
// Examples:
//
// Range expression ">=5.0.0 <6.0.0" contains 2 Constraints: ">=5.0.0" and "<6.0.0"
type Constraint struct {
	Operator     string
	Version      *Version
	IsXRange     bool
	IsTildeRange bool
	IsCaretRange bool
}

func (c *Constraint) String() string {
	return fmt.Sprintf("%s%s", c.Operator, c.Version.String())
}

// ParseRange from string
func ParseRange(s string) (*Range, error) {
	const (
		stateParseConstraint         = 1
		stateParseConstraintOperator = 2
		stateParseConstraintVersion  = 3
		stateParseConstraintOrSet    = 4
		stateParseSet                = 5
	)

	// special case "latest"
	if strings.ToLower(s) == "latest" {
		s = "*"
	}

	rang := &Range{Raw: s, Sets: make([]Set, 0)}
	set := make(Set, 0)
	constraint := &Constraint{}
	buffer := NewRuneBuffer()
	state := stateParseConstraint

	err := IterateString(s, func(r rune, iter *StringIteratorState) error {
		next, _ := iter.Next()

		switch state {
		case stateParseConstraint:
			newConstraint := false

			if isDigitRune(r) || r == 'x' || r == 'X' || r == '~' || r == '^' || r == '*' {
				newConstraint = true
				state = stateParseConstraintVersion
			}

			if isConstraintOperatorRune(r) {
				newConstraint = true
				state = stateParseConstraintOperator
			}

			if newConstraint {
				iter.Skip(-1)
				constraint = &Constraint{}
			}

			if !iter.HasNext() {
				return fmt.Errorf("Invalid range: %s", s)
			}

			return nil
		case stateParseConstraintOperator:
			if isConstraintOperatorRune(r) {
				buffer.Append(r)
			}

			if iter.HasNext() && isDigitRune(next) {
				constraint.Operator = buffer.String()
				buffer.Clear()
				state = stateParseConstraintVersion
				return nil
			}

			if !iter.HasNext() {
				return fmt.Errorf("Unterminated constraint: %s", s)
			}
		case stateParseConstraintVersion:
			if !isConstraintOperatorRune(r) && r != '|' {
				buffer.Append(r)
			}

			if r == '~' {
				if buffer.Len() != 1 {
					return fmt.Errorf("Invalid use of '~' in range: %s", s)
				}

				if constraint.Operator != "" {
					return fmt.Errorf(
						"Cannot use '~' in combination with operator '%s' in range: %s",
						constraint.Operator,
						s)
				}

				constraint.IsTildeRange = true
			}

			if r == '^' {
				if buffer.Len() != 1 {
					return fmt.Errorf("Invalid use of '^' in range: %s", s)
				}

				if constraint.Operator != "" {
					return fmt.Errorf(
						"Cannot use '^' in combination with operator '%s' in range: %s",
						constraint.Operator,
						s)
				}

				constraint.IsCaretRange = true
			}

			if isXRune(r) {
				// if constraint.Operator != "" {
				// 	return fmt.Errorf(
				// 		"Cannot use '%s' in combination with operator '%s' in range: %s",
				// 		string(r),
				// 		constraint.Operator,
				// 		s)
				// }

				constraint.IsXRange = true
				constraint.IsCaretRange = false
				constraint.IsTildeRange = false
			}

			if !iter.HasNext() || isConstraintOperatorRune(next) || next == '|' { //|| next == ' ' {
				if constraint.IsTildeRange {
					// Parse version containing ~ (ex: ~2.0.0)
					set, err := ParseTildeRangeSet(buffer.String())
					if err != nil {
						return err
					}
					rang.AddSet(set)
					state = stateParseConstraint

				} else if constraint.IsCaretRange {
					// Parse version containing ^ (ex: ^1.1.0)
					set, err := ParseCaretRangeSet(buffer.String())
					if err != nil {
						return err
					}
					rang.AddSet(set)
					state = stateParseConstraint

				} else if constraint.IsXRange {
					// Parse version containing x (ex: 1.x.x)
					if constraint.Operator != "" {
						// Example: ">=1.5.x"
						// in this scenario, we will just pad the version, replacing
						// x's with 0's, and parse this as ">=1.5.0"
						ver, err := ParseVersion(PadVersion(buffer.String()))
						if err != nil {
							return err
						}
						constraint.Version = ver
						set = append(set, constraint)
						if !iter.HasNext() {
							rang.AddSet(set)
							return nil
						}
					} else {
						set, err := ParseXRangeSet(buffer.String())
						if err != nil {
							return err
						}
						rang.AddSet(set)
						state = stateParseConstraint
					}

				} else if len(strings.Split(buffer.String(), ".")) < 3 {
					// Parse partial version (ex: 2.2)
					if constraint.Operator == "" {
						set, err := ParseXRangeSet(buffer.String())
						if err != nil {
							return err
						}
						rang.AddSet(set)
					} else {
						v := PadVersion(buffer.String())
						ver, err := ParseVersion(v)
						if err != nil {
							return err
						}
						constraint.Version = ver
						set = append(set, constraint)
						if !iter.HasNext() {
							rang.AddSet(set)
							return nil
						}
					}
					state = stateParseConstraint
				} else {
					// Default
					ver, err := ParseVersion(buffer.String())
					if err != nil {
						return err
					}
					constraint.Version = ver
					if constraint.Operator == "" {
						constraint.Operator = "="
					}
					set = append(set, constraint)
					if !iter.HasNext() {
						rang.AddSet(set)
						return nil
					}
				}

				constraint = &Constraint{}
				buffer.Clear()
				state = stateParseConstraintOrSet
			}

		case stateParseConstraintOrSet:
			if isDigitRune(r) || isConstraintOperatorRune(r) {
				state = stateParseConstraint
				iter.Skip(-1)
				return nil
			}

			if r == '|' {
				if len(set) > 0 {
					rang.AddSet(set)
				}
				buffer.Clear()
				state = stateParseSet
				iter.Skip(-1)
				return nil
			}
		case stateParseSet:
			if r == '|' {
				buffer.Append(r)
			}

			if iter.HasNext() && next != '|' {
				if buffer.String() != "||" {
					return fmt.Errorf("Unterminated logical OR operation in: %s", s)
				}
				set = make(Set, 0)
				buffer.Clear()
				state = stateParseConstraint
				return nil
			}

			if !iter.HasNext() {
				return fmt.Errorf("Invalid range: %s", s)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return rang, nil
}

// ParseTildeRangeSet parses a range containing a tilde
//
//    "~1.0.0" => ">=1.0.0 <1.1.0"
//
func ParseTildeRangeSet(s string) (Set, error) {
	if s[0] != '~' {
		return nil, fmt.Errorf("Invalid tilde range: %s", s)
	}

	// If there are missing parts (ex: ~1.1) then parse this as an X range
	if len(strings.Split(s, ".")) < 3 {
		s = PadVersion(s)
	}

	ver, err := ParseVersion(s[1:])
	if err != nil {
		return nil, err
	}

	set := make(Set, 0)
	set = append(set, &Constraint{Operator: ">=", Version: ver})
	set = append(set, &Constraint{Operator: "<", Version: &Version{Major: ver.Major, Minor: ver.Minor + 1, Patch: 0}})

	return set, nil
}

// ParseCaretRangeSet parses a range containing a caret
//
//    "^1.0.0" => ">=1.0.0 <2.0.0"
//
func ParseCaretRangeSet(s string) (Set, error) {
	if s[0] != '^' {
		return nil, fmt.Errorf("Invalid caret range: %s", s)
	}

	// If there are missing parts (ex: ~1.1) then parse this as an X range
	if len(strings.Split(s, ".")) < 3 {
		s = PadVersion(s)
	}

	ver, err := ParseVersion(s[1:])
	if err != nil {
		return nil, err
	}

	set := make(Set, 0)
	set = append(set, &Constraint{Operator: ">=", Version: ver})
	set = append(set, &Constraint{Operator: "<", Version: &Version{Major: ver.Major + 1, Minor: 0, Patch: 0}})

	return set, nil
}

// ParseXRangeSet parses a range containing an x or X
// Will also accept ranges with missing components (ex: 1.1)
//
//    "1.0.x" => ">=1.0.0 <1.1.0"
//    "1.8.x" => ">=1.8.0 <1.9.0"
//    "1.8" => ">=1.8.0 <1.9.0"
//    "1.x.x" => ">=1.0.0 <2.0.0"
//    "1.x" => ">=1.0.0 <2.0.0"
//    "x" => ">=0.0.0"
//
func ParseXRangeSet(s string) (Set, error) {
	if s[0] == '~' || s[0] == '^' {
		s = s[1:]
	}

	s = strings.TrimSpace(s)
	parts := strings.Split(s, ".")
	set := make(Set, 0)

	if len(parts) >= 1 && isXRune(rune(parts[0][0])) {
		set = append(set, &Constraint{Operator: ">=", Version: &Version{Major: 0, Minor: 0, Patch: 0}})
		return set, nil
	}

	majorInt, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("failed here, parsing: %s\n", s)
		return set, err
	}

	if len(parts) == 1 {
		set = append(set, &Constraint{Operator: ">=", Version: &Version{Major: majorInt, Minor: 0, Patch: 0}})
		set = append(set, &Constraint{Operator: "<", Version: &Version{Major: majorInt + 1, Minor: 0, Patch: 0}})
		return set, nil
	}

	if len(parts) >= 2 && isXRune(rune(parts[1][0])) {
		set = append(set, &Constraint{Operator: ">=", Version: &Version{Major: majorInt, Minor: 0, Patch: 0}})
		set = append(set, &Constraint{Operator: "<", Version: &Version{Major: majorInt + 1, Minor: 0, Patch: 0}})
		return set, nil
	}

	minorInt, err := strconv.Atoi(parts[1])
	if err != nil {
		return set, err
	}

	if len(parts) == 2 {
		set = append(set, &Constraint{Operator: ">=", Version: &Version{Major: majorInt, Minor: minorInt, Patch: 0}})
		set = append(set, &Constraint{Operator: "<", Version: &Version{Major: majorInt, Minor: minorInt + 1, Patch: 0}})
		return set, nil
	}

	if len(parts) >= 3 && isXRune(rune(parts[2][0])) {
		set = append(set, &Constraint{Operator: ">=", Version: &Version{Major: majorInt, Minor: minorInt, Patch: 0}})
		set = append(set, &Constraint{Operator: "<", Version: &Version{Major: majorInt, Minor: minorInt + 1, Patch: 0}})
		return set, nil
	}

	return set, fmt.Errorf("No wildcard (x, X, *) characters found in range: %s", s)
}
