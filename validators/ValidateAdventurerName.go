package validators

import "strings"

// The naming policy in BDO is fucked up
// This function will only check the length and allowed symbols
func ValidateAdventurerName(name *string) bool {
	if len(*name) < 3 || len(*name) > 16 {
		return false
	}

	f := func(r rune) bool {
		// Allowed characters: A-Z, a-z, 0-9, _
		return r < '0' || (r > '9' && r < 'A') || (r > 'Z' && r < '_') || (r > '_' && r < 'a') || r > 'z'
	}

	return strings.IndexFunc(*name, f) == -1
}
