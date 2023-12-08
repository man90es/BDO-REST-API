package validators

func ValidateRegionQueryParam(query []string) (region string) {
	if 1 > len(query) {
		return "EU"
	}

	// TODO: Add KR region once the translations are ready
	if query[0] == "NA" || query[0] == "SA" {
		return query[0]
	}

	return "EU"
}
