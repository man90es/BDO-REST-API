package validators

func ValidateSearchTypeQueryParam(query []string) string {
	if 1 > len(query) {
		return "2"
	}

	if query[0] == "characterName" {
		return "1"
	}

	return "2"
}
