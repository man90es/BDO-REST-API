package validators

import (
	"strings"
	"unicode"
)

// The naming policies in BDO are fucked up
// This function only checks the length and allowed symbols
func ValidateAdventurerNameQueryParam(query []string) (name string, ok bool) {
	if 1 > len(query) {
		return "", false
	}

	if len(query[0]) < 3 || len(query[0]) > 16 {
		return query[0], false
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

	return query[0], strings.IndexFunc(query[0], f) == -1
}
