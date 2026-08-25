package validators

// Check that the length is between 150 and 250 characters
func ValidateProfileTargetQueryParam(query []string) (profileTarget string, ok bool, errorMessage string) {
	if 1 > len(query) {
		return "", false, "Profile target is missing from the request"
	}

	if ok := len(query[0]) > 150 && len(query[0]) < 250; ok {
		return query[0], true, ""
	}

	return query[0], false, "Profile target has to be between 150 and 250 characters long"
}
