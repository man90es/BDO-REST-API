package validators

import (
	"fmt"
	"strings"
	"unicode"
)

// The naming policies in BDO are fucked up
// This function only checks the length and allowed symbols
// searchType: 1 for character, 2 for family
func ValidateAdventurerNameQueryParam(query []string, region string, searchType string) (name string, ok bool, errorMessage string) {
	if 1 > len(query) {
		return "", false, "Adventurer name is missing from request"
	}

	// Length bounds have been found for each region by trial and error,
	// they're wider than the game says it allows
	minLength := map[string]int{
		"EU1": 3,
		"EU2": 1,
		"KR1": 2,
		"KR2": 2,
		"NA1": 3,
		"NA2": 1,
		"SA1": 2,
		"SA2": 2,
	}[region+searchType]
	maxLength := map[string]int{
		"EU1": 19,
		"EU2": 16,
		"KR1": 10,
		"KR2": 10,
		"NA1": 28,
		"NA2": 24,
		"SA1": 18,
		"SA2": 16,
	}[region+searchType]

	name = strings.ToLower(query[0])

	if len(name) < minLength {
		return name, false, fmt.Sprintf("Adventurer name on %v should be at least %v symbols long", region, minLength)
	}

	if len(name) > maxLength {
		return name, false, fmt.Sprintf("Adventurer name on %v should be at most %v symbols long", region, maxLength)
	}

	// Returns false for allowed characters
	// and true for everything else
	f := func(r rune) bool {
		// Numbers
		if unicode.IsNumber(r) {
			return false
		}

		// Latin letters
		if unicode.Is(unicode.Latin, r) {
			return false
		}

		// Underscore
		if r == '_' {
			return false
		}

		// Korean characters
		if unicode.Is(unicode.Hangul, r) {
			return false
		}

		return true
	}

	if i := strings.IndexFunc(name, f); i != -1 {
		return name, false, fmt.Sprintf("Adventurer name contains a forbidden symbol at position %v: %q", i+1, query[0][i])
	}

	return name, true, ""
}
