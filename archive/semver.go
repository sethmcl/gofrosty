package lib

import (
	"fmt"
	"github.com/sethmcl/gofrosty/vendor/semver"
	"regexp"
	"strconv"
	"strings"
)

// ConvertNpmSemver converts an npm semver to an expression supported by the
// github.com/blang/semver library. Can be compound expression.
//
// Examples:
//    "~1.0.0" returns ">=1.0.0 <1.1.0"
//    "2 || 3 || 4" returns ">=2.0.0 <3.0.0 || >=3.0.0 <4.0.0 || >=4.0.0 <5.0.0"
func ConvertNpmSemver(v string) (string, error) {
	compound := strings.Split(v, "||")
	if len(compound) == 1 {
		return ConvertSingleNpmSemver(v)
	}

	var parts []string
	for _, part := range compound {
		converted, err := ConvertSingleNpmSemver(part)
		if err != nil {
			return "", err
		}
		parts = append(parts, converted)
	}
	return strings.Join(parts, " || "), nil
}

// ConvertSingleNpmSemver converts an npm semver to an expression supported by the
// github.com/blang/semver library
//
// Examples:
//    "~1.0.0" returns ">=1.0.0 <1.1.0"
//    "1.0.x" returns ">=1.0.0 <1.1.0"
//    "^1.3.0" returns ">=1.3.0 <2.0.0"
//    "2.3.0" returns "2.3.0"
//    "2 || 3" returns ">=3.0.0 <4.0.0"
func ConvertSingleNpmSemver(v string) (string, error) {
	type Ver struct {
		Major       string
		Minor       string
		Patch       string
		latestMinor bool
		latestPatch bool
	}

	GetContext().Info("parsing %s", v)
	// v = PadSemver(v)
	parts := strings.Split(v, ".")
	vObj := &Ver{"0", "0", "0", false, false}
	if len(parts) >= 1 {
		vObj.Major = parts[0]
	}
	if len(parts) >= 2 {
		vObj.Minor = parts[1]
	}
	if len(parts) >= 3 {
		vObj.Patch = parts[2]
	}

	// If patch not specified, use latest patch
	if len(parts) == 2 {
		vObj.latestPatch = true
	}

	// If minor not specified, use latest minor
	if len(parts) == 1 {
		vObj.latestMinor = true
	}

	// If * is first character, then use latest version
	if v[0] == '*' {
		return ">=0.0.0", nil
	}

	// If ^ is first character, then use latest minor version
	if v[0] == '^' {
		vObj.latestMinor = true
	}

	// If ~ is first character, then use lastest patch
	if v[0] == '~' {
		vObj.latestPatch = true
	}

	stripRe, err := regexp.Compile("[~|^|v|V|=| |>|<]")
	if err != nil {
		return "", err
	}
	vObj.Major = stripRe.ReplaceAllString(vObj.Major, "")

	if strings.ToUpper(string(vObj.Minor[0])) == "X" {
		vObj.latestMinor = true
		vObj.Minor = "0"
	}

	if strings.ToUpper(string(vObj.Patch[0])) == "X" {
		vObj.latestPatch = true
		vObj.Patch = "0"
	}

	// Strip extensions and convert to ints
	stripExtensionRe, err := regexp.Compile("[-|+].*$")
	if err != nil {
		return "", err
	}

	majorStrip := stripExtensionRe.ReplaceAllString(vObj.Major, "")
	majorInt, err := strconv.Atoi(majorStrip)
	if err != nil {
		return "", err
	}

	minorStrip := stripExtensionRe.ReplaceAllString(vObj.Minor, "")
	minorInt, err := strconv.Atoi(minorStrip)
	if err != nil {
		return "", err
	}

	patchStrip := stripExtensionRe.ReplaceAllString(vObj.Patch, "")
	patchInt, err := strconv.Atoi(patchStrip)
	if err != nil {
		return "", err
	}

	// Handle ranges
	if vObj.latestMinor {
		return fmt.Sprintf(
			">=%d.%d.%d <%d.0.0",
			majorInt,
			minorInt,
			patchInt,
			majorInt+1), nil
	}

	if vObj.latestPatch {
		return fmt.Sprintf(
			">=%d.%d.%d <%d.%d.0",
			majorInt,
			minorInt,
			patchInt,
			majorInt,
			minorInt+1), nil

	}

	return strings.Join(parts, "."), nil
}

// IsExplicitSemver returns true if string is an explicit semver
func IsExplicitSemver(v string) bool {
	result, err := ConvertNpmSemver(v)
	if err != nil {
		return false
	}

	spaces, err := regexp.MatchString(" ", v)
	if err != nil {
		return false
	}

	specialChar, err := regexp.MatchString("[<|>|=|\\|]", v)
	if err != nil {
		return false
	}

	return len(strings.Split(result, ".")) == 3 && !spaces && !specialChar
}

// IsSemverRange returns true if string is a semver range
func IsSemverRange(v string) bool {
	canParseConvertedRange := false
	canParseRange := false
	isExplicit := IsExplicitSemver(v)

	_, err := semver.ParseRange(v)
	if err == nil {
		canParseRange = true
	}

	cv, err := ConvertNpmSemver(v)
	if err == nil {
		_, err = semver.ParseRange(cv)
		if err == nil {
			canParseConvertedRange = true
		} else {
			GetContext().Debug("Error calling semver.ParseRange(%s)", cv)
		}
	}

	return !isExplicit && (canParseRange || canParseConvertedRange)
}

// ResolveSemver resolves to newest semver which matches semver.Range conditions
func ResolveSemver(v string, candidates []string) (string, error) {
	selectedCandidate := ""
	npmv, err := ConvertNpmSemver(v)
	if err == nil {
		v = npmv
	}

	svCandidates, err := SemverVersionsFromStrings(candidates)
	if err != nil {
		return "", err
	}

	srange, err := semver.ParseRange(v)
	if err != nil {
		return "", err
	}

	idx := len(svCandidates) - 1
	for idx >= 0 {
		candidate := svCandidates[idx]
		if srange(candidate) {
			selectedCandidate = candidate.String()
			break
		}
		idx--
	}

	if selectedCandidate == "" {
		return "", fmt.Errorf("Cannot resolve %s", v)
	}

	return selectedCandidate, nil
}

// SemverVersionsFromStrings converts semver's expressed as strings to semver.Version objects
func SemverVersionsFromStrings(vs []string) ([]semver.Version, error) {
	var result []semver.Version
	for _, v := range vs {
		sver, err := semver.Parse(v)
		if err != nil {
			return nil, err
		}

		result = append(result, sver)
	}
	semver.Sort(result)
	return result, nil
}

// PadSemver add missing parts to semver string
func PadSemver(v string) string {
	re, err := regexp.Compile("([0-9|.]+)")
	if err != nil {
		return v
	}

	sv := re.ReplaceAllStringFunc(v, func(match string) string {
		parts := strings.Split(match, ".")
		if len(parts) < 3 {
			for idx := 0; idx < 3-len(parts); idx++ {
				match += ".0"
			}
		}
		return match
	})

	symRe, err := regexp.Compile("([>|=|<])[ ]+([0-9])")
	return symRe.ReplaceAllString(sv, "$1$2")
}

// // ParseSemver
