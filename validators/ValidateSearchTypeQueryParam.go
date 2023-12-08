package validators

func ValidateSearchTypeQueryParam(query []string) (searchType string) {
	if 1 > len(query) {
		return "familyName"
	}

	if query[0] == "characterName" {
		return query[0]
	}

	return "familyName"
}
