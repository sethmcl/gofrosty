package nsemver

var rangeRunes = []rune{'<', '>', '|', 'x', 'X', '~', '^'}
var digitRunes = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var constraintOperatorRunes = []rune{'<', '>', '='}
var versionSeparatorRunes = []rune{'.', '-'}
var setSeparatorRunes = []rune{'|'}
var constraintSeparatorRunes = []rune{'-'}
var xRunes = []rune{'x', 'X', '*'}

// sliceContainsRune returns true if rune is found within rune slice
func sliceContainsRune(haystack []rune, needle rune) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}

// isRangeRune returns true if rune character indicates a range string
func isRangeRune(r rune) bool {
	return sliceContainsRune(rangeRunes, r)
}

// isDigitRune returns true if rune character is a digit
func isDigitRune(r rune) bool {
	return sliceContainsRune(digitRunes, r)
}

// isVersionSeparatorRune returns true if rune character is a separator
func isVersionSeparatorRune(r rune) bool {
	return sliceContainsRune(versionSeparatorRunes, r)
}

// isSetSeparatorRune returns true if rune character is a separator
func isSetSeparatorRune(r rune) bool {
	return sliceContainsRune(setSeparatorRunes, r)
}

// isConstraintSeparatorRune returns true if rune character is a separator
func isConstraintSeparatorRune(r rune) bool {
	return sliceContainsRune(constraintSeparatorRunes, r)
}

// isConstraintOperatorRune returns true if rune character is a constraint operator
func isConstraintOperatorRune(r rune) bool {
	return sliceContainsRune(constraintOperatorRunes, r)
}

// isXRune returns true if rune character is an x rune
func isXRune(r rune) bool {
	return sliceContainsRune(xRunes, r)
}
