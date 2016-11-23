package lib

// SliceContains returns true if slice contains value
func SliceContains(haystack []interface{}, needle interface{}) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

// StringSliceContains returns true if slice contains value
func StringSliceContains(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}
