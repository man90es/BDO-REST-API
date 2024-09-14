package validators

import (
	"fmt"
	"strings"
	"unicode"
)

// The naming policies in BDO are fucked up
// This function only checks the length and allowed symbols
func ValidateAdventurerNameQueryParam(query []string, region string) (name string, ok bool, errorMessage string) {
	if 1 > len(query) {
		return "", false, "Adventurer name is missing from request"
	}

	minLength := map[string]int{
		"SA": 2,
		"NA": 3,
		"EU": 3,
	}[region]

	if len(query[0]) < minLength || len(query[0]) > 16 {
		return query[0], false, fmt.Sprintf("Adventurer name should be between %v and 16 symbols long", minLength)
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

	if i := strings.IndexFunc(query[0], f); i != -1 {
		return query[0], false, fmt.Sprintf("Adventurer name contains a forbidden symbol at position %v: %q", i+1, query[0][i])
	}

	return query[0], true, ""
}
