package validators

// Check that the length is at least 150 characters
// I don't actually know how long it should be, but the length varies
func ValidateProfileTargetQueryParam(query []string) (profileTarget string, ok bool, errorMessage string) {
	if 1 > len(query) {
		return "", false, "Profile target is missing from the request"
	}

	if ok := len(query[0]) >= 150; ok {
		return query[0], true, ""
	}

	return query[0], false, "Profile target has to be at least 150 characters long"
}
