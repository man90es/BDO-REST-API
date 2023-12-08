package validators

import "strings"

func ValidateRegionQueryParam(query []string) (region string) {
	if 1 > len(query) {
		return "EU"
	}

	region = strings.ToUpper(query[0])

	// TODO: Add KR region once the translations are ready
	if region == "NA" || region == "SA" {
		return region
	}

	return "EU"
}
