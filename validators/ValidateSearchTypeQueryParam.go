package validators

func ValidateSearchTypeQueryParam(query []string) (searchType uint8, searchTypeAsString string) {
	if 1 > len(query) {
		return 2, "familyName"
	}

	if query[0] == "characterName" {
		return 1, "characterName"
	}

	return 2, "familyName"
}
