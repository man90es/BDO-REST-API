package validators

// Check that the length is at least 150 characters
// I don't actually know how long it should be, but the length varies
func ValidateProfileTarget(profileTarget *string) bool {
	return len(*profileTarget) >= 150
}
