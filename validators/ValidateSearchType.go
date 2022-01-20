package validators

func ValidateSearchType(s *string) bool {
	return *s == "characterName" || *s == "familyName"
}
