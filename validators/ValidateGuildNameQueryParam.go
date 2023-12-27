package validators

import (
	"strings"
	"unicode"
)

// The naming policies in BDO are fucked up
// This function only checks the length and allowed symbols
// I also assumed that the allowed symbols are the same as for adventurer names
func ValidateGuildNameQueryParam(query []string) (guildName string, ok bool) {
	if 1 > len(query) {
		return "", false
	}

	guildName = strings.ToLower(query[0])

	if len(guildName) < 2 {
		return guildName, false
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

	return guildName, strings.IndexFunc(guildName, f) == -1
}
