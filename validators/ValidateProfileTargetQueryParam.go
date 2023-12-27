package validators

// Check that the length is at least 150 characters
// I don't actually know how long it should be, but the length varies
func ValidateProfileTargetQueryParam(query []string) (profileTarget string, ok bool) {
	if 1 > len(query) {
		return "", false
	}

	return query[0], len(query[0]) >= 150
}
