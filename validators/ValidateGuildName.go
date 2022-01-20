package validators

import "strings"

func ValidateGuildName(name *string) bool {
	if len(*name) < 2 {
		return false
	}

	f := func(r rune) bool {
		// Implying that allowed characters are the same as with player names
		return r < '0' || (r > '9' && r < 'A') || (r > 'Z' && r < '_') || (r > '_' && r < 'a') || r > 'z'
	}

	return strings.IndexFunc(*name, f) == -1
}
