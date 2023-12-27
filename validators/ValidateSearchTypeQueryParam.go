package validators

func ValidateSearchTypeQueryParam(query []string) (searchType uint8) {
	if 1 > len(query) {
		return 2
	}

	if query[0] == "characterName" {
		return map[string]uint8{
			"characterName": 1,
			"familyName":    2,
		}[query[0]]
	}

	return 2
}
